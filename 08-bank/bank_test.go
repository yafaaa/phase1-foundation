package main

import "testing"

func TestBankAccount(t *testing.T) {
	tests := []struct {
		name        string
		initial     int
		action      string // "deposit" or "withdraw"
		amount      int
		wantBalance int
		wantErr     bool
	}{
		{"Valid deposit", 100, "deposit", 50, 150, false},
		{"Negative deposit", 100, "deposit", -30, 100, true},
		{"Valid withdraw", 100, "withdraw", 40, 60, false},
		{"Overdraft withdraw", 100, "withdraw", 150, 100, true},
		{"Negative withdraw", 100, "withdraw", -10, 100, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			acc := &BankAccount{balance: tc.initial}
			var err error

			if tc.action == "deposit" {
				err = acc.Deposit(tc.amount)
			} else {
				err = acc.Withdraw(tc.amount)
			}

			if (err != nil) != tc.wantErr {
				t.Errorf("got err = %v, wantErr = %v", err, tc.wantErr)
			}
			if acc.Balance() != tc.wantBalance {
				t.Errorf("got balance = %d, want = %d", acc.Balance(), tc.wantBalance)
			}
		})
	}
}