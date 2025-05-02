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
	Screen := myApp.NewWindow("Hello")

	name_Text := widget.NewLabel("")

	mainUiInnerWindows := container.NewHBox()

	newRoomNameEntry := widget.NewEntry()

	chats := container.NewVBox()

	sideBar := container.NewHSplit(
		container.NewVBox(
			name_Text,

			chats,

			widget.NewButton("New Chat", func() {
				mainUiInnerWindows.Add(container.NewInnerWindow("pls",
					container.NewVBox(
						container.NewGridWithColumns(
							2,
							widget.NewLabel("Chat Name"),
							newRoomNameEntry,
						),
						widget.NewButton("Create Chat", func() {
							data := NetworkChat{
								newRoomNameEntry.Text,
								[]string{Current_User.Username},
							}
							temp_data, err := json.Marshal(data)
							if err != nil {
								panic(err)
							}
							_, err = http.Post("http://localhost:5151/MakeChat", "json", bytes.NewBuffer(temp_data))
							if err != nil {
								panic(err)
							}

							chats.Add(widget.NewButton(data.Name, func() {
							}))
						}),
					)))
			}),
		),
		container.NewVBox(),
	)
	sideBar.Offset = -0.1

	mainUi := container.NewVBox(
		sideBar,
		mainUiInnerWindows,
	)

	loginUsernameInput := widget.NewEntry()
	loginUsernameInput.SetPlaceHolder("")

	loginUsernameContainer := container.NewGridWithColumns(
		2,
		widget.NewLabel("Username: "),
		loginUsernameInput,
	)

	loginPasswordInput := widget.NewEntry()
	loginPasswordInput.SetPlaceHolder("")

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

			Screen.SetContent(mainUi)

			Current_User = NewUser(loginUsernameInput.Text, loginPasswordInput.Text)
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
				name_Text.Text += Current_User.Username

				Screen.SetContent(mainUi)

			}

			getRoomData, err := json.Marshal(NetworkChat{Current_User.Username, []string{}})

			resp, err = http.Post("http://localhost:5151/GetChatsForUser", "json", bytes.NewBuffer(getRoomData))
			if err != nil {
				panic(err)
			}

			temp_string_chats, err := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}

			chats_data := []Chat{}
			if err := json.Unmarshal(temp_string_chats, &chats_data); err != nil {
				panic(err)
			}

			for _, chat := range chats_data {
				chats.Add(widget.NewButton(chat.Name, func() {
					// to do later
				}))
			}
		}),
	)

	Screen.SetContent(
		loginUi,
	)

	Screen.Show()
	myApp.Run()
	tidyUp()
}

func tidyUp() {
	fmt.Println("Exited")
}
