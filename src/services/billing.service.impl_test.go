package services

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"testing"

// 	"github.com/stripe/stripe-go/v81"
// )

// func TestBillingServiceStripeImpl_CreateOrder(t *testing.T) {
// 	fmt.Printf("\"aa.\": %v\n", ".")
// 	type fields struct {
// 		appSuccessUrl string
// 		appCancelUrl  string
// 	}
// 	type args struct {
// 		ctx          context.Context
// 		currencyUnit stripe.Currency
// 		unitAmmount  int64
// 		planName     string
// 	}
// 	tests := []struct {
// 		name    string
// 		fields  fields
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "one",
// 			fields: fields{
// 				"http://127.0.0.1/success",
// 				"http://127.0.0.1/cancel",
// 			},
// 			args: args{
// 				ctx:          context.Background(),
// 				currencyUnit: stripe.CurrencyBRL,
// 				unitAmmount:  1.99 * 100,
// 				planName:     "QuackWeek! - Event Registration",
// 			},
// 			wantErr: false,
// 		},
// 	}
// 	stripe.Key = os.Getenv("STRIPE_API_KEY")
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			s := &BillingServiceStripeImpl{
// 				appSuccessUrl: tt.fields.appSuccessUrl,
// 				appCancelUrl:  tt.fields.appCancelUrl,
// 			}
// 			got, err := s.CreateOrder(tt.args.ctx, tt.args.currencyUnit, tt.args.unitAmmount, tt.args.planName, 1)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("BillingServiceStripeImpl.CreateOrder() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			fmt.Printf("got: %v\n", got)
// 		})
// 	}
// }

// func TestBillingServiceStripeImpl_GetCheckoutSession(t *testing.T) {
// 	t.Run("a", func(t *testing.T) {
// 		stripe.Key = os.Getenv("STRIPE_API_KEY")
// 		s := &BillingServiceStripeImpl{
// 			db:            nil,
// 			appSuccessUrl: "tt.fields.appSuccessUrl",
// 			appCancelUrl:  "tt.fields.appCancelUrl",
// 		}
// 		got, err := s.GetCheckoutSession(context.Background(), "SESSION_ID")
// 		if err != nil {
// 			t.Errorf("BillingServiceStripeImpl.GetCheckoutSession() error = %v", err)
// 			return
// 		}

// 		fmt.Printf("got: %v\n", got)
// 	})
// }
