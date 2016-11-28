package main

import (
    "fmt"
    "strings"
    "io/ioutil"

    "github.com/go-gl/gl/v3.2-core/gl"
)

type Shader struct {
    program uint32
}

func (s *Shader) Use() {
    gl.UseProgram(s.program)
}

func compile(source string, shaderType uint32) (uint32, error) {
    shader := gl.CreateShader(shaderType)

    csources, free := gl.Strs(source)
    gl.ShaderSource(shader, 1, csources, nil)
    free()
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("Failed to compile %v: %v", source, log)
    }

    return shader, nil
}

func NewShader(vertex, fragment string) (*Shader, error) {
    s := new (Shader)

    vertexShader, err := compile(vertex, gl.VERTEX_SHADER)
    if err != nil {
        return nil, err
    }

    fragmentShader, err := compile(fragment, gl.FRAGMENT_SHADER)
    if err != nil {
        return nil, err
    }

    s.program = gl.CreateProgram()
    gl.AttachShader(s.program, vertexShader)
    gl.AttachShader(s.program, fragmentShader)
    gl.LinkProgram(s.program)

    var status int32
    gl.GetProgramiv(s.program, gl.LINK_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetProgramiv(s.program, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetProgramInfoLog(s.program, logLength, nil, gl.Str(log))

        return nil, fmt.Errorf("Failed to ling program: %v", log)
    }

    gl.DeleteShader(vertexShader)
    gl.DeleteShader(fragmentShader)

    return s, nil
}

func LoadShader(vertexFile, fragmentFile string) (*Shader, error) {
    // Load shader files
    vertexBuf, err := ioutil.ReadFile("vertex.glsl")
    if err != nil {
        return nil, err
    }
    vertex := string(vertexBuf)

    fragmentBuf, err := ioutil.ReadFile("fragment.glsl")
    if err != nil {
        return nil, err
    }
    fragment := string(fragmentBuf)

    // Compile shader
    shader, err := NewShader(vertex, fragment)
    if err != nil {
        return nil, err
    }

    return shader, nil
}
