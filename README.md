# AcountLandBot

An exclusively made Telegram Bot for AcountLand ([Aryaat](https://t.me/Aryaaat))

This bot is able to generate the images split into multiple tabs of Fortnite's Item Shop customized with their preferences.

> The rewritten version of this bot is built using the [tgo](https://github.com/haashemi/tgo) framework; a ⭐️ would be appreciated.

![20231123_172918](https://github.com/haashemi/AcountLandBot/assets/60406325/92e23c70-adf6-4909-b791-3e5bdac64beb)

## Commands:

- `/itemshop` Generates the tabs' images and sends them to you one by one.
- `/setpp` It updates the primary price rate and rewrites the `config.yaml` file for you.
- `/setsp` It updates the secondary price rate and rewrites the `config.yaml` file for you.

## Usage:

1. Copy or rename `config.example.yaml` to `config.yaml` and modify with your requirements.
2. Update the background image in `generator/assets/images/background.png` if you want to. (you should)
3. Build the app using `go build .`
4. Run the app! (using `./AcountLandBot`)

## Acknowledges:

- There's an issue with some Farsi letters for primary and secondary titles; in that case, it is recommended to use Arabic letters if you can.

  - It's caused by the lack of support for RTL letters, and the upstream package [garabic](github.com/abdullahdiaa/garabic) doesn't handle some Farsi letters.

- If you want to use an English title, you'll need to do some code or asset modifications.
  - Font files are at `generator/assets/fonts`
  - Fonts are loaded at `generator/loaders.go` in `loadFonts`
  - Fonts and Titles are generated [here](https://github.com/haashemi/AcountLandBot/blob/3034aa7b4ff77e01f86bc8fa45ee940e7e4db5bd/generator/itemshop.go#L101-L110)

- It only generate images with Outfits and bundles; you can modify it [here](https://github.com/haashemi/AcountLandBot/blob/458265bcaa61d102e778c3bd0ab26c9661b0b661/bot/itemshop.go#L107-L119)

## Clarification:

This project was initially private-source (December 2021–November 2023), but as priorities changed and time passed, I and the customer (MR. Arya) talked a little and decided to open-source this project after a rewrite. And here is the result!

## Thanks to:

- [Fortnite-API.com](https://fortnite-api.com) for their awesome API.
