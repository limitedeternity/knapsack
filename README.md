# knapsack
> Просто рюкзаки. Разные: ограниченные и неограниченные

## Cборка:

1. Нужен Go 1.22

2. Просто `go build`, и всё

## Использование:

```powershell
❯ .\knapsack.exe -h
Options:

  -h, --help                 Display help information
  -i, --items               *Yaml file with an array of items
  -c, --capacity            *Knapsack capacity
  -k, --knapsack[=bounded]   Knapsack type
```

* Через `-i` указывается путь до .yaml-файла со списком предметов, которые нужно попробовать впихнуть в рюкзак (см. [items.yaml](https://github.com/limitedeternity/knapsack/blob/master/items.yaml))

* Через `-c` задаётся ёмкость рюкзака, в которую надо вписаться

* Через `-k` можно указать, какой тип рюкзака использовать: ограниченный (`bounded`; по-умолчанию) или неограниченный (`unbounded`)

## Схема items.yaml:

* `item`: `string` (required; название предмета)

* `weight`: `integer` (required; вес предмета)

* `value` : `integer` (required; ценность предмета)

* `pieces`: `integer` (optional; количество единиц предмета)

Поле `pieces` нужно для ограниченного рюкзака, неограниченному на него всё равно.
Если его не указать для ограниченного рюкзака, то поле примет значение по-умолчанию (`1`), и будет решаться задача про рюкзак 0/1.

## Пример:

```powershell
❯ .\knapsack.exe -i items.yaml -c 8
Taking:
+ 2m: 1
+ 6m: 1
Total value: 22
Total weight: 8
```

Да, это Rod Cutting, самый первый тест-кейс.
