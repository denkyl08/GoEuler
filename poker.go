package main

import (
	"fmt"
	"os"
	"bufio"
	"strings"
	"time"
)

type Card struct {
	value, suit byte 
}

type Hand []Card


/*func (hand Hand) ContainsCard(value, suit byte) bool {
	for i:=0; i < len(hand); i++ {
		if hand[i].value== value && hand[i].suit == suit {
			return true
		}
	}
	return false
}*/

func (hand Hand) Value() (value int, cardValue byte) {
	/*
	0 - High Card: Highest value card.
	1 - One Pair: Two cards of the same value.
	2 - Two Pairs: Two different pairs.
	3 - Three of a Kind: Three cards of the same value.
	4 - Straight: All cards are consecutive values.
	5 - Flush: All cards of the same suit.
	6 - Full House: Three of a kind and a pair.
	7 - Four of a Kind: Four cards of the same value.
	8 - Straight Flush: All cards are consecutive values of same suit.
	9 - Royal Flush: Ten, Jack, Queen, King, Ace, in same suit.
	*/

	//Sort the hand by value
	for card:=0; card<len(hand); card++ {
		for i:=0; i<len(hand)-1; i++ {
			if hand[i].value > hand[i+1].value {
				heavyCard := hand[i+1]
				hand[i+1] = hand[i]
				hand[i] = heavyCard
			}
		}
	}
	
	var Straight bool = true
	for i:=0; i < len(hand)-1; i++ {
		if hand[i].value + byte(1) != hand[i+1].value { Straight=false; break }
	}

	var Flush bool = false
	if  hand[0].suit == hand[1].suit && 
		hand[0].suit == hand[2].suit && 
		hand[0].suit == hand[3].suit && 
		hand[0].suit == hand[4].suit {
		Flush = true	
	}
	if Flush &&  Straight && hand[0].value == byte(8) { //8 is index of 'T' in order array
	   return 9, hand[4].value    //Royal Flush
	}

	if Straight && Flush {
		return 8, hand[4].value   //Straight Flush
	}

	var kindsCountArray []int
	var kindsArray []Card
	for i:=0; i < len(hand); i++ {
		found:=false
		for kind:=0; kind< len(kindsArray); kind++{
			if hand[i].value == kindsArray[kind].value { 
				kindsCountArray[kind] += 1; found= true; break 
			}
		}
		if found == false || len(kindsArray)==0 {
			kindsCountArray = append(kindsCountArray, 1)
			kindsArray = append(kindsArray, hand[i])
		}
	}

	var FourOfAKind bool = false
	var ThreeOfAKind bool = false
	var TwoOfAKind bool = false
	var HighestKindIndex int
	for i:=0; i < len(kindsCountArray); i++ {
		if kindsCountArray[i] == 4 {FourOfAKind=true;HighestKindIndex=i;break}
		if kindsCountArray[i] == 3 {ThreeOfAKind=true; HighestKindIndex=i}
		if kindsCountArray[i] == 2 {
			TwoOfAKind=true
			if !ThreeOfAKind{HighestKindIndex=i}
		}
	}

	if FourOfAKind {
		return 7, kindsArray[HighestKindIndex].value      //Four of a kind
	}
	if TwoOfAKind && ThreeOfAKind {
		return 6, kindsArray[HighestKindIndex].value      //Full House
	}
	if Straight {
		return 5, hand[4].value    //Straight
	}
	if Flush {
		return 4, hand[4].value     //Flush
	}
	if ThreeOfAKind {
		return 3, kindsArray[HighestKindIndex].value      //Three of a kind
	}
	if TwoOfAKind {
		if len(kindsArray)<=3 {
			return 2, kindsArray[HighestKindIndex].value  //Two Pair
		} else {
			return 1, kindsArray[HighestKindIndex].value  //One Pair
		}
		
	}
	return 0, hand[4].value
}

func StringToByteHand(str string) (hand Hand) {
	OrderedByValue := []byte{'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A'}
	StrArray := strings.Split(str, " ")
	for i:= 0; i<len(StrArray); i++ {
		var card Card;
		for index:=0; index < len(OrderedByValue); index++ {
			if []byte(StrArray[i])[0] == OrderedByValue[index] {card.value = byte(index);break}
		}
		card.suit = []byte(StrArray[i])[1]	
		hand = append(hand, card)
	}
	return hand
}

func ByteHandToString(hand Hand) (str string) {
	OrderedByValue := []byte{'2', '3', '4', '5', '6', '7', '8', '9', 'T', 'J', 'Q', 'K', 'A'}
	for i:=0; i<len(hand); i++ {
		str = str + string(OrderedByValue[hand[i].value]) +string(hand[i].suit)
		if i != len(hand)-1 {str = str + " "}
	}
	return str
}


func HighestHand(players []Hand, winCount []int) (winner int){
	HighestHand := -1
	HighCard := byte(0)
	for hand:=0; hand<len(players); hand++ {
		CurrentHand, ThisHighCard := players[hand].Value()
		if  CurrentHand > HighestHand {
			HighestHand = CurrentHand
			winner = hand
			HighCard = ThisHighCard
		} else if CurrentHand == HighestHand {
			if ThisHighCard>HighCard {
				winner = hand
				HighCard = ThisHighCard
			}
		}

	}
	winCount[winner] = winCount[winner]+1
	return winner
}


func main() {
	beginning:=time.Now()
	fmt.Println("Begin Euler Problem #54: Poker Hands")
	file, err := os.Open("poker.txt"); if err != nil {
		return
	}
	reader := bufio.NewReader(file)
	winCount := []int{0, 0}
	for {
		str, err := reader.ReadString('\n')
		hand := StringToByteHand(str)
		var players []Hand
		for i:=0; i+5 <= 10; i+=5 {
			players = append(players, hand[i:i+5])
		} 
		go HighestHand(players, winCount)
		if err != nil { break }
	}
	fmt.Println("player 1 won ", winCount[0], "times")
	fmt.Println("completed in", time.Since(beginning).Nanoseconds()/1000/1000, "milliseconds")
}