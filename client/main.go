package main

import (
	"context"
	"time"

	pb "github.com/SzymonMielecki/chatApp/usersService"
	"github.com/rivo/tview"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Screen string

const (
	unloggedScreen Screen = "unlogged"
	mainScreen     Screen = "main"
	loginScreen    Screen = "login"
	registerScreen Screen = "register"
)

type App struct {
	client        pb.UsersServiceClient
	app           *tview.Application
	currentScreen Screen
	screens       map[Screen]tview.Primitive
}

func NewApp(c pb.UsersServiceClient) *App {
	app := &App{
		client:        c,
		app:           tview.NewApplication(),
		currentScreen: unloggedScreen,
		screens:       make(map[Screen]tview.Primitive),
	}

	app.screens[unloggedScreen] = app.createUnloggedScreen()
	app.screens[loginScreen] = app.createLoginScreen()
	app.screens[registerScreen] = app.createRegisterScreen()
	app.screens[mainScreen] = app.createMainScreen()

	return app
}

func (app *App) Run(initialScreen Screen) {
	app.SwitchToScreen(initialScreen)
	if err := app.app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) SwitchToScreen(screen Screen) {
	app.currentScreen = screen
	app.app.SetRoot(app.screens[screen], true).SetFocus(app.screens[screen])
}

func (app *App) createUnloggedScreen() tview.Primitive {
	unlogged := tview.NewList()
	unlogged.AddItem("Login", "if you already have an account", 'l', func() {
		app.SwitchToScreen(loginScreen)
	})
	unlogged.AddItem("Register", "if you don't have an account", 'r', func() {
		app.SwitchToScreen(registerScreen)
	})
	unlogged.AddItem("Quit", "Press to exit", 'q', func() {
		app.app.Stop()
	})
	return unlogged
}

func (app *App) createLoginScreen() tview.Primitive {
	username := ""
	password := ""
	loginForm := tview.NewForm()
	loginForm.AddInputField("Username", "", 20, nil, func(text string) {
		username = text
	})
	loginForm.AddPasswordField("Password", "", 20, '*', func(text string) {
		password = text
	})
	loginForm.AddButton("Login", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		r, err := app.client.LoginUser(ctx, &pb.LoginUserRequest{UsernameOrEmail: username, Password: password})
		if err != nil {
			app.SwitchToScreen(unloggedScreen)
		}
		if r.Success {
			app.SwitchToScreen(mainScreen)
		} else {
			app.SwitchToScreen(unloggedScreen)
		}
	})
	loginForm.AddButton("Back", func() {
		app.SwitchToScreen(unloggedScreen)
	})
	return loginForm
}

func (app *App) createRegisterScreen() tview.Primitive {
	username := ""
	password := ""
	confirmPassword := ""
	registerForm := tview.NewForm()
	registerForm.AddInputField("Username", "", 20, nil, func(text string) {
		username = text
	})
	registerForm.AddPasswordField("Password", "", 20, '*', func(text string) {
		password = text
	})
	registerForm.AddPasswordField("Confirm Password", "", 20, '*', func(text string) {
		confirmPassword = text
	})
	registerForm.AddButton("Register", func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		if password != confirmPassword {
			app.app.Stop()
		}
		r, err := app.client.RegisterUser(ctx, &pb.RegisterUserRequest{Username: username, Password: password})
		if err != nil {
			app.SwitchToScreen(unloggedScreen)
		}
		if r.Success {
			app.SwitchToScreen(mainScreen)
		} else {
			app.SwitchToScreen(unloggedScreen)
		}
	})
	registerForm.AddButton("Back", func() {
		app.SwitchToScreen(unloggedScreen)
	})
	return registerForm
}

func (app *App) createMainScreen() tview.Primitive {
	// Implement your main screen here
	mainScreen := tview.NewTextView().SetText("Welcome to the main screen!")
	return mainScreen
}

func main() {
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	c := pb.NewUsersServiceClient(conn)

	app := NewApp(c)
	app.Run(unloggedScreen)
}
