package gpt

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// Offer represents the structure for an offer.
type Offer struct {
	Cashback     float32 `json:"cashback"`
	Condition    string  `json:"condition"`
	Expiry       string  `json:"expiry"`
	Restrictions string  `json:"restrictions"`
	Category     string  `json:"category"`
	CardType     string
	BankID       int
}

type Offers struct {
	Offers []Offer `json:"offers"`
}

// GPTClient struct to hold any configuration
type GPTClient struct {
	Client *openai.Client
	ctx    context.Context
}

func NewClient(api string, ctx context.Context) GPTClient {
	client := openai.NewClient(api)
	return GPTClient{client, ctx}
}

func (client GPTClient) SendRequest(prompt string) (string, error) {
	resp, err := client.Client.CreateChatCompletion(
		client.ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}
	return resp.Choices[0].Message.Content, nil
}

// AnalyzeOffers sends a request to GPT to analyze cashback offers based on the bank, card type, and data provided.
func (client GPTClient) AnalyzeOffers(bank int, cardType, parsedData string) ([]Offer, error) {
	// Correctly escaped prompt
	prompt := fmt.Sprintf(`Analyze these text: %s. Return me answer in json format. I will parse answer like array of: type Offer struct {
		
		Cashback     float32 json:"cashback" //it is cashback precent in float You should extract it from data.If you cant find just write random int between 3 and 10.round to one decimal place
		Category    string  json:"category" // it is name of offer.
		Expiry       string  json:"expiry" //expiration date. null by default, you can put random integers from 1 to 100
		Condition    string  json:"condition"
		Restrictions string  json:"restrictions" //add "only weekdays" or "only weekends"
		

	}. like:{
        "Offers": [
                {
                    "cashback": 15.0,
                    "category": Gaming Services,
                    "expiry": null,
					"condition": More than 3000
                    "restrictions": null,
                    
                }
				{
					"cashback": 12.0,
                    "category": Beauty Saloon,
                    "expiry": null,
					"condition": More than 5000
                    "restrictions": null,
                        
				}
        ]
}Lenght of array should be 14 in total. on for each of these Category:'Gaming Services', 'Beauty Saloon', 'Clothes and Shoes', 'Furniture', 'Medical Services', 'Cafes and Restaurants', 'Taxi', 'Online Cinema and Music', 'Travel', 'Fitness and SPA', 'Supermarket', 'Education', 'Food Delivery', 'Products for Children'`, parsedData)

	response, err := client.SendRequest(prompt)
	if err != nil {
		return nil, err
	}
	//fmt.Println(response)

	var offers Offers
	err = json.Unmarshal([]byte(response), &offers)
	if err != nil {
		return nil, errors.New("failed to parse JSON response: " + err.Error())
	}

	return offers.Offers, nil
}
