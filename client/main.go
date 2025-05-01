package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var Current_User = User{}

func main() {
	myApp := app.New()
	loginScreen := myApp.NewWindow("Hello")

	loginUsernameInput := widget.NewEntry()
	loginUsernameInput.SetPlaceHolder("                ")

	loginUsernameContainer := container.NewGridWithColumns(2,
		widget.NewLabel("Username: "),
		loginUsernameInput,
	)

	loginPasswordInput := widget.NewEntry()
	loginPasswordInput.SetPlaceHolder("                ")

	loginPasswordContainer := container.NewGridWithColumns(2,
		widget.NewLabel("Password: "),
		loginPasswordInput,
	)

	loginPasswordContainer.Resize(fyne.NewSize(100, 100))
	loginUsernameContainer.Resize(fyne.NewSize(100, 100))

	loginUi := container.NewVBox(
		widget.NewLabel("Login"),
		loginUsernameContainer,
		loginPasswordContainer,
		widget.NewButton("Create User", func() {
			data, err := json.Marshal(NewUser(loginUsernameInput.Text, loginPasswordInput.Text))
			if err != nil {
				panic(err)
			}
			resp, err := http.Post("http://localhost:5151/AddUser", "json", bytes.NewBuffer(data))
			if err != nil {
				panic(err)
			}
			defer resp.Body.Close()

			fmt.Println(resp.Status)
		}),

		widget.NewButton("Login", func() {
			data, err := json.Marshal(NewUser(loginUsernameInput.Text, loginPasswordInput.Text))
			if err != nil {
				panic(err)
			}
			resp, err := http.Post("http://localhost:5151/GetUser", "json", bytes.NewBuffer(data))
			if err != nil {
				panic(err)
			}

			if resp.StatusCode == 202 {
				temp_data, err := io.ReadAll(resp.Body)
				if err != nil {
					panic(err)
				}

				temp_user := User{}

				json.Unmarshal(temp_data, &temp_user)

				Current_User = temp_user
			}

			fmt.Println(Current_User)
		}),
	)

	loginScreen.SetContent(
		loginUi,
	)

	loginScreen.Show()
	myApp.Run()
	tidyUp()
}

func tidyUp() {
	fmt.Println("Exited")
}
