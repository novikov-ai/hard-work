# Задачка N 0

## Что делает код

~~~python
def fibs_sum(): # генератор 1
    fsum = 0
    while True:
        fsum += yield fsum
     
def get_fibs(num): # генератор 2
    a, b = 0, 1
    gsum = fibs_sum()

    gsum.send(None) # инициализация генератора пустым значением
    for i in range(num): # итерации по передаваемому количеству
        yield gsum.send(b) # отправка нового состояния в генератор и возвращение состояния накопительным итогом
        c = b
        b = a + b
        a = c
~~~

## Исправленный

~~~python   
def get_fibs(num):
    a, b = 0, 1
    
    for i in range(num):
        yield b
        a, b = b, a+b
~~~