## API

### Инициализация
#### Создать игру
##### Request
* `URL: [POST] /xo`
* `Body: {name: ":name"}`
##### Response
[Объект комнаты (Room)](#Room)

### Получить список комнат
#### Request
* `URL: [GET] /room`
#### Response
[Список объектов комнат (Room)](#Room)

### Создать нового игрока в комнате
#### Request
* `URL: [POST] /xo/:roomUID/player`
  * `:roomUID` uid из объекта [Room](#Room)
* `Body: {name: ":name"}`
#### Response
[Объект игрока](#Player)
#### Errors
* `room with UID: ':roomUID' doesn't exists` - комната не существует

### Создать нового наблюдателя
#### Request
* `URL: [POST] /xo/:roomUID/watcher`
  * `:roomUID` uid из объекта [Room](#Room)
* `Body: {name: ":name"}`
#### Response
[Объект наблюдателя](#Watcher)
#### Errors
* `room with UID: ':roomUID' doesn't exists` - комната не существует

### Websocket-ручка
#### Connect
* `URL: [WS] /xo/:roomUID/connect/:participantUID`
    * `:roomUID` - `:roomUID` uid из объекта [Room](#Room)
    * `:participantUID` - `:participantUID` uid из объекта [Participant](#Participant)
#### User Events
##### Сходить
```json
{
    "type": "xo_move", 
    "x": 0, 
    "y": 0
}
```
`x, y` - координаты поля на доске
#### Server Events
##### Состояние доски
```json
{
  "type": "xo_board_state",
  "payload": [
    ["-", "-", "-"],
    ["-", "-", "-"],
    ["-", "-", "-"]
  ]
}
```
Элементы `payload` может принимать значения: `"x"`, `"o"`, `"-"`
##### Ошибка хода
```json
{
  "type": "xo_move_error",
  "payload": ":sting"
}
``` 
Возможные значения `payload`:
* `"field is occupied"` - клетка занята
* `"game is finish"` - игра окончена
* `"waiting opponent moving"` - ожидаение хода оппонента
##### Окончание игры для наблюдателя
```json
{
  "type": "xo_watcher_game_result",
  "payload": ":sting"
}
```
`payload` указывает на победившего игрока. Возможные значения `payload`: `"x"`, `"o"`, `"-"`
##### Окончание игры для игрока
```json
{
  "type": "xo_player_game_result",
  "payload": ":sting"
}
```
`payload` сообщает результат относительно игрока. Возможные значения `payload`: `"xo_win"`, `"xo_loose"`, `"xo_draw"`
#### Errors
* `room with UID: ':roomUID' doesn't exists` - комната не существует
* `participant with UID: ':participantUID' doesn't exists"` - участник не существует
* `participant already connected` - слот участника занят

### Объекты
#### Room
```json
{
    "uid": ":uuid",
    "name": ":string",
    "participants": ":array<Participant>",
    "meta": ":<RoomMeta>"
}
```

#### Participant
```json
{
    "uid": ":uuid",
    "name": ":string",
    "connected": ":bool",
    "meta": ":<ParticipantMeta>"
}
```

#### XO Room Meta (RoomMeta) 
```json
{
    "type": ":string"
}
```
Константа: `{"type": "xo"}`

#### Player 
Participant Meta (ParticipantMeta)
```json
":string"
```
* `"x-player"` для игрока за "крестики"
* `"o-player"` для игрока за "нолики"

#### Watcher
Participant Meta (ParticipantMeta)
```json
":string"
```
Константа `"watcher"`

#### Error
Ошибки API
```json
{"message": ":string"}
```
