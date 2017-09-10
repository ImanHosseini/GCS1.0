// Copyright 2014 The go-gl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Renders a textured spinning cube using GLFW 3 and OpenGL 4.1 core forward-compatible profile.
package main // import "github.com/go-gl/example/gl41core-cube"

import (
	"fmt"
	"go/build"
	"image"
	"image/draw"
	_ "image/png"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"

	//"math"

	_ "math/rand"
	_ "math"
	"math"
)

const windowWidth = 1400
const windowHeight = 900
const N = 40
const yzero = 1.0
const dt=0.01
const g=0.5
const damp=0.9
const k=9.0*float32(N)
const wind_amp=5.0
const wind_tau=100.0
const dx = 2.0/float32(N)

type Particle struct{
	pos mgl32.Vec3
	vel mgl32.Vec3
	acc mgl32.Vec3
}

type World struct{
	nodes []Particle
	timer int
}

func InitWorld() (*World){
	var nodes []Particle
	nodes = make([]Particle,(N+1)*(N+1))
	for i:=0; i<=N; i++{
		for j:=0; j<=N; j++{
			x:= -1.0+float32(i)*dx
			z:= -1.0+float32(j)*dx
			var ind = i*(N+1)+j
			//Vertex
			nodes[ind]=Particle{mgl32.Vec3{x,yzero,z},mgl32.Vec3{0.0,0.0,0.0},mgl32.Vec3{0.0,0.0,0.0}}
		}
	}
	var w World=World{nodes,0}
	return &w
}

func calcForce(p1,p2 Particle) (mgl32.Vec3) {
	delta := p2.pos.Sub(p1.pos)
	return  delta.Mul(k*(1.0-(delta.Len()/dx)))
}

func (w *World) update(){
	//fmt.Println(w.timer)
	(*w).timer++
	for i:=N+1; i<len(w.nodes); i++ {
		v:=&((w).nodes[i])
		(*v).acc = mgl32.Vec3{0.0,0.0,0.0}
		if i%(N+1)!=0 {
			(*v).acc=v.acc.Add(calcForce(w.nodes[i-1],*v))
		}
		if (i+1)%(N+1)!=0 {
			(*v).acc=v.acc.Add(calcForce(w.nodes[i+1],*v))
		}
		if (i+N+1)<(N+1)*(N+1) {
			(*v).acc=v.acc.Add(calcForce(w.nodes[i+1+N],*v))
		}
		if (i-N-1)>-1 {
			(*v).acc=v.acc.Add(calcForce(w.nodes[i-1-N],*v))
		}
		(*v).acc = (v).acc.Add(mgl32.Vec3{0.0,-g,0.0})
		(*v).acc=(v).acc.Sub(v.vel.Mul(damp))
		//if w.timer>750{
		//	var pow = (1.5+0.5*rand.Float32())
		//	var dist2 =float32 (math.Hypot (float64(v.pos.Y()+0.5+0.3*(rand.Float32()-0.5)),float64(v.pos.Z()+0.3*(rand.Float32()-0.5))))
		//	dist2 *= 3.0
		//	if dist2>1.0{
		//		pow = pow/(dist2*dist2)
		//	}
		//
		//	//if rand.Float32()>0.999{
		//	//	conz=-conz
		//	//}
		//	(*v).acc=v.acc.Add(mgl32.Vec3{-pow*conz,0.0,0.0})
		//}
	}
	for i:=N+1; i<len(w.nodes); i++ {
		v := &((w).nodes[i])
		(*v).vel= v.vel.Add(v.acc.Mul(dt))
		(*v).pos= v.pos.Add(v.vel.Mul(dt))
		//(*v).pos=mgl32.Vec3{0.0,11110.0,0.0}
	}

	var i1beg = int (N/8)
	if(w.timer>750 && w.timer%250==0){
		for i:=i1beg; i<7*i1beg; i++{
			for j:=i1beg; j<7*i1beg; j++{
				var pow = 2.0
				var softner = math.Hypot(float64(i-N/2),float64(j-N/2))
				softner *= 3.0*16.0/(7.0*float64(N))
				if softner>1.0 {
					pow /= (softner*softner)
				}
				ind := i*(N+1)+j
				if w.timer%500==250 {
					pow= (-pow)
				}
				v:=&((w).nodes[ind])
				(*v).vel=v.vel.Add(mgl32.Vec3{float32 (pow),0.0,0.0})
			}
		}
		//for i:=5*i1beg; i<7*i1beg; i++{
		//	for j:=i1beg; j<7*i1beg; j++{
		//		ind := i*(N+1)+j
		//		v:=&((w).nodes[ind])
		//		(*v).vel=v.vel.Add(mgl32.Vec3{-1.5,0.0,0.0})
		//	}
		//}
	}
	//if(w.timer==1500){
	//	for i:=5; i<10; i++{
	//		for j:=5; j<10; j++{
	//			ind := i*(N+1)+j
	//			v:=&((w).nodes[ind])
	//			(*v).vel=v.vel.Add(mgl32.Vec3{3.0,0.0,0.0})
	//		}
	//	}
	//}

}

func (w World) drawV(){
	for i:=0; i<N; i++{
		for j:=0; j<N; j++{
			var base = (i*N+j)*30
			var baseind = i*(N+1)+j
			//Vertex - (X Y Z U V)
			vertices2[base] = w.nodes[baseind].pos.X()
			vertices2[base+1] = w.nodes[baseind].pos.Y()
			vertices2[base+2] = w.nodes[baseind].pos.Z()
			vertices2[base+3] = 0.0
			vertices2[base+4] = 0.0
			//Vertex - (X Y Z U V)
			vertices2[base+5] = w.nodes[baseind+1].pos.X()
			vertices2[base+6] = w.nodes[baseind+1].pos.Y()
			vertices2[base+7] = w.nodes[baseind+1].pos.Z()
			vertices2[base+8] = 1.0
			vertices2[base+9] = 0.0
			//Vertex - (X Y Z U V)
			vertices2[base+10] = w.nodes[baseind+N+1].pos.X()
			vertices2[base+11] = w.nodes[baseind+N+1].pos.Y()
			vertices2[base+12] = w.nodes[baseind+N+1].pos.Z()
			vertices2[base+13] = 0.0
			vertices2[base+14] = 1.0
			//Vertex - (X Y Z U V)
			vertices2[base+15] = w.nodes[baseind+1].pos.X()
			vertices2[base+16] = w.nodes[baseind+1].pos.Y()
			vertices2[base+17] = w.nodes[baseind+1].pos.Z()
			vertices2[base+18] = 1.0
			vertices2[base+19] = 0.0
			//Vertex - (X Y Z U V)
			vertices2[base+20] = w.nodes[baseind+N+2].pos.X()
			vertices2[base+21] = w.nodes[baseind+N+2].pos.Y()
			vertices2[base+22] = w.nodes[baseind+N+2].pos.Z()
			vertices2[base+23] = 1.0
			vertices2[base+24] = 1.0
			//Vertex - (X Y Z U V)
			vertices2[base+25] =  w.nodes[baseind+N+1].pos.X()
			vertices2[base+26] = w.nodes[baseind+N+1].pos.Y()
			vertices2[base+27] =  w.nodes[baseind+N+1].pos.Z()
			vertices2[base+28] = 0.0
			vertices2[base+29] = 1.0
		}
	}
}


var vertices2 []float32
var w World
var conz bool

func init() {
	conz=true
    vertices2 = make([]float32,N*N*30)
	w = *InitWorld()
	(w).drawV()
	// GLFW event handling must run on the main OS thread
	runtime.LockOSThread()
}


func main() {
	if err := glfw.Init(); err != nil {
		log.Fatalln("failed to initialize glfw:", err)
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	window, err := glfw.CreateWindow(windowWidth, windowHeight, "Cube", nil, nil)
	if err != nil {
		panic(err)
	}
	window.MakeContextCurrent()

	// Initialize Glow
	if err := gl.Init(); err != nil {
		panic(err)
	}

	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version", version)

	// Configure the vertex and fragment shaders
	program, err := newProgram(vertexShader, fragmentShader)
	if err != nil {
		panic(err)
	}

	gl.UseProgram(program)

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(windowWidth)/windowHeight, 0.1, 10.0)
	projectionUniform := gl.GetUniformLocation(program, gl.Str("projection\x00"))
	gl.UniformMatrix4fv(projectionUniform, 1, false, &projection[0])

	camera := mgl32.LookAtV(mgl32.Vec3{3, 3, 3}, mgl32.Vec3{0, 0, 0}, mgl32.Vec3{0, 1, 0})
	cameraUniform := gl.GetUniformLocation(program, gl.Str("camera\x00"))
	gl.UniformMatrix4fv(cameraUniform, 1, false, &camera[0])

	model := mgl32.Ident4()
	modelUniform := gl.GetUniformLocation(program, gl.Str("model\x00"))
	gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

	textureUniform := gl.GetUniformLocation(program, gl.Str("tex\x00"))
	gl.Uniform1i(textureUniform, 0)

	gl.BindFragDataLocation(program, 0, gl.Str("outputColor\x00"))

	// Load the texture
	texture, err := newTexture("square.png")
	if err != nil {
		log.Fatalln(err)
	}

	// Configure the vertex data
	var vao uint32
	gl.GenVertexArrays(1, &vao)
	gl.BindVertexArray(vao)

	var vbo uint32
	gl.GenBuffers(1, &vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices2)*4, gl.Ptr(vertices2), gl.STATIC_DRAW)

	vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
	gl.EnableVertexAttribArray(vertAttrib)
	gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

	texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
	gl.EnableVertexAttribArray(texCoordAttrib)
	gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.ClearColor(1.0, 1.0, 1.0, 1.0)

	angle := 0.0
	previousTime := glfw.GetTime()

	for !window.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

		// Update
		time := glfw.GetTime()
		elapsed := time - previousTime
		previousTime = time

		angle += elapsed*0.0
		//fmt.Println(elapsed)
		model = mgl32.HomogRotate3D(float32(angle), mgl32.Vec3{0, 1, 0})

		(w).update()
		(w).drawV()
		//fmt.Println(w.nodes[50].pos)
		// End Update

		// Render Update
		gl.UseProgram(program)
		gl.UniformMatrix4fv(modelUniform, 1, false, &model[0])

		gl.BindVertexArray(vao)

		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, texture)
		//
		// Configure the vertex data
		var vao uint32
		gl.GenVertexArrays(1, &vao)
		gl.BindVertexArray(vao)

		var vbo uint32
		gl.GenBuffers(1, &vbo)
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
		gl.BufferData(gl.ARRAY_BUFFER, len(vertices2)*4, gl.Ptr(vertices2), gl.STATIC_DRAW)

		vertAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vert\x00")))
		gl.EnableVertexAttribArray(vertAttrib)
		gl.VertexAttribPointer(vertAttrib, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))

		texCoordAttrib := uint32(gl.GetAttribLocation(program, gl.Str("vertTexCoord\x00")))
		gl.EnableVertexAttribArray(texCoordAttrib)
		gl.VertexAttribPointer(texCoordAttrib, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))

		// Configure global settings
		gl.Enable(gl.DEPTH_TEST)
		gl.DepthFunc(gl.LESS)
		gl.ClearColor(1.0, 1.0, 1.0, 1.0)
		// Render Update
		gl.DrawArrays(gl.TRIANGLES, 0, N*N*2*3)

		// Maintenance
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

func newProgram(vertexShaderSource, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()

	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)

		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(log))

		return 0, fmt.Errorf("failed to link program: %v", log)
	}

	defer gl.DeleteShader(vertexShader)
	defer gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
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

		return 0, fmt.Errorf("failed to compile %v: %v", source, log)
	}

	return shader, nil
}

func newTexture(file string) (uint32, error) {
	imgFile, err := os.Open(file)
	if err != nil {
		return 0, fmt.Errorf("texture %q not found on disk: %v", file, err)
	}
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, fmt.Errorf("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}

var vertexShader = `
#version 330

uniform mat4 projection;
uniform mat4 camera;
uniform mat4 model;

in vec3 vert;
in vec2 vertTexCoord;

out vec2 fragTexCoord;

void main() {
    fragTexCoord = vertTexCoord;
    gl_Position = projection * camera * model * vec4(vert, 1);
}
` + "\x00"

var fragmentShader = `
#version 330

uniform sampler2D tex;

in vec2 fragTexCoord;

out vec4 outputColor;

void main() {
    outputColor = texture(tex, fragTexCoord);
}
` + "\x00"

var cubeVertices = []float32{
	//  X, Y, Z, U, V
	// Bottom
	-1.0, -1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,

	// Top
	-1.0, 1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, 1.0, 0.0, 1.0,
	1.0, 1.0, 1.0, 1.0, 1.0,

	// Front
	-1.0, -1.0, 1.0, 1.0, 0.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, 1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,

	// Back
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	-1.0, 1.0, -1.0, 0.0, 1.0,
	1.0, 1.0, -1.0, 1.0, 1.0,

	// Left
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,
	-1.0, -1.0, -1.0, 0.0, 0.0,
	-1.0, -1.0, 1.0, 0.0, 1.0,
	-1.0, 1.0, 1.0, 1.0, 1.0,
	-1.0, 1.0, -1.0, 1.0, 0.0,

	// Right
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, -1.0, -1.0, 1.0, 0.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, -1.0, 1.0, 1.0, 1.0,
	1.0, 1.0, -1.0, 0.0, 0.0,
	1.0, 1.0, 1.0, 0.0, 1.0,
}

// Set the working directory to the root of Go package, so that its assets can be accessed.
func init() {
	dir, err := importPathToDir("github.com/go-gl/example/gl41core-cube")
	if err != nil {
		log.Fatalln("Unable to find Go package in your GOPATH, it's needed to load assets:", err)
	}
	err = os.Chdir(dir)
	if err != nil {
		log.Panicln("os.Chdir:", err)
	}
}

// importPathToDir resolves the absolute path from importPath.
// There doesn't need to be a valid Go package inside that import path,
// but the directory must exist.
func importPathToDir(importPath string) (string, error) {
	p, err := build.Import(importPath, "", build.FindOnly)
	if err != nil {
		return "", err
	}
	return p.Dir, nil
}