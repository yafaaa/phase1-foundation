package main

import "errors"

type BankAccount struct {
	balance int
}

func (b *BankAccount) Deposit(amount int) error {
	// TODO: If amount <= 0, return errors.New("deposit must be positive")
	// Otherwise add amount to b.balance and return nil
	if amount <= 0{
		return errors.New("Deposite must be positive")
	}

	b.balance += amount
	return nil
}

func (b *BankAccount) Withdraw(amount int) error {
	// TODO:
	
	if amount <= 0{
		return errors.New("withdrew amount must be positive")
	}else if amount > b.balance{
		return errors.New("insufficient funds / overdraft")
	}else{
		b.balance -= amount
		return nil
	}
}

func (b *BankAccount) Balance() int {
	// TODO: return b.balance
	return b.balance
}