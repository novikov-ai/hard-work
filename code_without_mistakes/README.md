# Пишем безошибочный код

Есть класс `Acount`, который позволял выполнять различные операции: пополнение и снятие денег, закрытие и открытие счета. Операции могли вызываться в произвольном порядке, что требовало дополнительных проверок в каждом методе, чтобы избежать некорректных состояний. 

До:
~~~go
type Account struct {
    balance int
    isClosed bool
}

func (a *Account) Deposit(amount int) {
    if a.isClosed {
        panic("Cannot deposit to a closed account")
    }
    a.balance += amount
}

func (a *Account) Withdraw(amount int) {
    if a.isClosed {
        panic("Cannot withdraw from a closed account")
    }
    a.balance -= amount
}

func (a *Account) Close() {
    a.isClosed = true
}

func (a *Account) Open() {
    a.isClosed = false
}
~~~

После рефакторинга общий класс `Acount` был разбит на несколько: `OpenAccount` и `ClosedAccount`, которые позволяли совершать определенные операции, произвольный порядок которых не мог в принципе нарушить работу систему или привести к некорректному состоянию.

После:
~~~go
type OpenAccount struct {
    balance int
}

func (o *OpenAccount) Deposit(amount int) {
    o.balance += amount
}

func (o *OpenAccount) Withdraw(amount int) {
    o.balance -= amount
}

func (o *OpenAccount) Close() ClosedAccount{
     // implementation
   return ClosedAccount{}
}

type ClosedAccount struct {
    balance int
}

func (c *ClosedAccount) Open() OpenAccount{
    // implementation
   return OpenAccount{}
}
~~~