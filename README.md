# Pong
One day at night, I was watching [Johnathan Blow's Video](https://youtu.be/nL8GWU9M8LY?si=fdT7x8pIXU4AJ_EG), and I heard him saying, "Stop doing web programming". 😔 we all know we do web cause we kinda have to. 
So, I decided the next day to make a simple thing, a "Pong" game with Raylib. So I made a one-person Pong game, and then I thought, "What good is pong if you don't play it with your friends?" So I tried to make it playable by LAN, and here it is.

## Prerequisite

- [go-raylib](https://github.com/gen2brain/raylib-go) and it's prerequisites.

## How to run?
Make sure you are on the same LAN ( WLAN ) as your friend. 
and make sure you install all the dependencies by running the following command.

```bash
go mod tidy
```
### To host a game
Run the following
```bash
go run main.go --host
```
and then it will log your <SERVER_IP>
### To join
run the following command by replacing `<SERVER_IP>` with your logged server IP. ( don't include the port )
```bash
go run main.go --join <SERVER_IP>
```

## How to play?
You will use HJKL to move your pads. If you know Vim motion, you know what I mean. but if you need further.

```
             ^
             k             
       < h       l >               
             j                    
             v
```

