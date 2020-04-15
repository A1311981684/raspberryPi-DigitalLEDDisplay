package LED3461BS

import (
	"fmt"
	rpi "github.com/nathan-osman/go-rpigpio"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"
)

//显示从 0.000 到 9999 的数字
//Display Number from 0.000 to 9999

//使用树莓派4B实现四位LED数码管的基本控制[使用BCM引脚编码]
//Implement 4-Number Digital LED Display in Raspberry Pi 4B with Linux OS [PIN coded with BCM]

//默认引脚连接方式(我兴许大概是画对了，不然可能我吨吨吨有点多了)
//Default PIN connection(This should be right, I checked it...unless I am drunk):
/***********************************************************************************
##简介##
控制PIN 1、2、3、4用于控制操作哪个数字（小数点处于下方时，从左往右依次标记为1、2、3、4）。
将控制PIN设置为高电平则表示它会被操作并显示。

显示PIN A,B,C,D,E,F,G用于控制7段线段组成一个0~9的数字（或者其他你认为有意义的图案）。

###Brief Introduce###:
Control PORT 1,2,3,4 are used to control target Single Number(With dots at the bottom,
 from left to right, they are indexed as 1,2,3,4).
Set the Control PORT to "High" means it is enabled to display.

Display Port A,B,C,D,E,F,G are used to control the 7 sticks in a Digital Number to form
 a number within 0~9(or something else means something to you).

GPIO:       0    1   27    2    3    4
BCM:       17   18   16   27    22   23
PORT:		1	 A	  F	   2	3	 B
			|	 |	  |	   |	|	 |
          -----   -----   -----   -----
         |     | |     | |     | |     |
         |     | |     | |     | |     |
	      -----   -----   -----   -----
         |     | |     | |     | |     |
         |     | |     | |     | |     |
	      ----- . ----- . ----- . -----    (Notice: there are those Dots(Control by DP Pin) int this line)
		   |	 |	  |	   |	|	 |
PORT:      E	 D	  DP   C    G    4
GPIO:      26    6    29   5    28   21
 BCM:      12    25   21   24   20   5

*************************************************************************************
单个数字控制方式：
Single number control:

       (A)
      -----
  (F)|     |(B)
 	 |     |
	  -(G)-
  (E)|     |(C)
 	 |     |
	  -----  .(DP)
       (D)

************************************************************************************/

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
}

//be creative, add more chars
var (
	//order: A B C D E F G
	//zero means light, one means go off
	num0  = []int{0, 0, 0, 0, 0, 0, 1}
	num1  = []int{1, 0, 0, 1, 1, 1, 1}
	num2  = []int{0, 0, 1, 0, 0, 1, 0}
	num3  = []int{0, 0, 0, 0, 1, 1, 0}
	num4  = []int{1, 0, 0, 1, 1, 0, 0}
	num5  = []int{0, 1, 0, 0, 1, 0, 0}
	num6  = []int{0, 1, 0, 0, 0, 0, 0}
	num7  = []int{0, 0, 0, 1, 1, 1, 1}
	num8  = []int{0, 0, 0, 0, 0, 0, 0}
	num9  = []int{0, 0, 0, 0, 1, 0, 0}
	charH = []int{1, 0, 0, 1, 0, 0, 0}
	charC = []int{0, 1, 1, 0, 0, 0, 1}
	//a 2D slice, with index as its meaning, except for non-number content, I'm saying, "H", "C"
	nums = [][]int{num0, num1, num2, num3, num4, num5, num6, num7, num8, num9, charH, charC}
)

type Led3461BS struct {
	controlPin1, controlPin2, controlPin3, controlPin4,
	displayPinA, displayPinB, displayPinC, displayPinD, displayPinE, displayPinF, displayPinG,
	dpPin *rpi.Pin

	initAlready bool
	finished    chan bool
}

//Init PINs
func (l *Led3461BS) Init() (err error) {
	l.controlPin1, err = rpi.OpenPin(17, rpi.OUT)
	if err != nil {
		return
	}
	l.controlPin2, err = rpi.OpenPin(27, rpi.OUT)
	if err != nil {
		return
	}
	l.controlPin3, err = rpi.OpenPin(22, rpi.OUT)
	if err != nil {
		return
	}
	l.controlPin4, err = rpi.OpenPin(5, rpi.OUT)
	if err != nil {
		return
	}
	l.displayPinA, err = rpi.OpenPin(18, rpi.OUT)
	if err != nil {
		return
	}
	l.displayPinB, err = rpi.OpenPin(23, rpi.OUT)
	if err != nil {
		return
	}
	l.displayPinC, err = rpi.OpenPin(24, rpi.OUT)
	if err != nil {
		return
	}
	l.displayPinD, err = rpi.OpenPin(25, rpi.OUT)
	if err != nil {
		return
	}
	l.displayPinE, err = rpi.OpenPin(12, rpi.OUT)
	if err != nil {
		return
	}
	l.displayPinF, err = rpi.OpenPin(16, rpi.OUT)
	if err != nil {
		return
	}
	l.displayPinG, err = rpi.OpenPin(20, rpi.OUT)
	if err != nil {
		return
	}
	l.dpPin, err = rpi.OpenPin(21, rpi.OUT)
	if err != nil {
		return
	}
	l.initAlready = true
	l.finished = make(chan bool, 0)

	return nil
}

//Ignored now
func (l *Led3461BS) Release() error {
	//Reserved
	err := l.controlPin1.Close()
	if err != nil {
		return err
	}
	err=l.controlPin2.Close()
	if err != nil {
		return err
	}
	err=l.controlPin3.Close()
	if err != nil {
		return err
	}
	err=l.controlPin4.Close()
	if err != nil {
		return err
	}
	err=l.dpPin.Close()
	if err != nil {
		return err
	}
	err=l.displayPinA.Close()
	if err != nil {
		return err
	}
	err=l.displayPinB.Close()
	if err != nil {
		return err
	}
	err=l.displayPinC.Close()
	if err != nil {
		return err
	}
	err=l.displayPinD.Close()
	if err != nil {
		return err
	}
	err=l.displayPinE.Close()
	if err != nil {
		return err
	}
	err=l.displayPinF.Close()
	if err != nil {
		return err
	}
	err=l.displayPinG.Close()
	if err != nil {
		return err
	}
	return nil
}

//Display requested content
func (l *Led3461BS) Execute(content interface{}, durationSec int) error {

	timer := time.NewTimer(time.Duration(durationSec) * time.Second)
	c := make(chan int, 0)
	go func() {
		for ;; {
			select {
			case <-c:
				return
			default:
			}
			err := l.process(content)
			if err != nil {
				log.Println(err.Error())
				_ = l.Dark()

				break
			}
		}
	}()

	//Display over
	select {
	case <-timer.C:
		_ = l.Dark()
		c<-1
	}

	return nil
}

//start displaying content(a number)
func (l *Led3461BS) process(content interface{}) error {

	err := l.Dark()
	if err != nil {
		return err
	}
	//must init first
	if !(l.initAlready) {
		return fmt.Errorf("PINs are not initialized, please invoke this.Init() first")
	}

	//string assertion
	displayNum, ok := content.(string)
	if !ok {
		return fmt.Errorf("invalid param: value:%v, value:%v", reflect.ValueOf(content), reflect.TypeOf(content))
	}

	//char count must less than 6 and positive
	charNumber := len([]byte(displayNum))
	if charNumber >= 6 || charNumber == 0 {
		return fmt.Errorf("error input: chars length is not valid:%s, length:%d, want 0<LENGTH<6", displayNum, len([]byte(displayNum)))
	}

	//valid content must be able to convert into a float
	converted, err := strconv.ParseFloat(displayNum, 32)
	if err != nil {
		return err
	}
	if converted < 0.0 || converted > 9999.0 {
		return fmt.Errorf("number out of range:%.4f, need 0.0 <= 9999.0 ", converted)
	}
	//analyze content
	//first get all single char
	var temp = strings.Replace(displayNum, ".", "", -1)
	//11.1
	indexOfDot := strings.LastIndex(displayNum, ".")
	dotNo := 0
	if indexOfDot >= 0 {
		dotNo = 4 - (len([]byte(displayNum)[indexOfDot+1:]))
	} else {
		dotNo = -1
	}
	//I like to light up dot first!
	err = l.LightUpDot(dotNo)
	if err != nil {
		return err
	}
	//write 4 chars quickly
	//With the help of outside loop, the frequency of LED refreshing won't be noticed by human eyes
	for i := 0; i < len([]byte(temp)); i++ {
		start := 4 - len([]byte(temp)) + 1 + i
		intV, err := strconv.Atoi(string(temp[i]))
		if err != nil {
			return err
		}
		err = l.DisplaySingleChar(start, intV)
		if err != nil {
			return err
		}
	}
	return nil
}

//choose one char to display at a specific index
func (l *Led3461BS) DisplaySingleChar(no, value int) error {
	err := l.Dark()
	if err != nil {
		return err
	}
	err = l.controlPin1.Write(rpi.Value(0))
	if err != nil {
		return err
	}
	err = l.controlPin2.Write(rpi.Value(0))
	if err != nil {
		return err
	}
	err = l.controlPin3.Write(rpi.Value(0))
	if err != nil {
		return err
	}
	err = l.controlPin4.Write(rpi.Value(0))
	if err != nil {
		return err
	}
	switch no {
	case 1:
		err := l.controlPin1.Write(rpi.HIGH)
		if err != nil {
			return err
		}
	case 2:
		err := l.controlPin2.Write(rpi.HIGH)
		if err != nil {
			return err
		}
	case 3:
		err := l.controlPin3.Write(rpi.HIGH)
		if err != nil {
			return err
		}
	case 4:
		err := l.controlPin4.Write(rpi.HIGH)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("invalid no for control PIN:%d", no)
	}
	return l.pinWrite(nums[value])
}

//write Display PIN by pre-defined slice
func (l *Led3461BS) pinWrite(input []int) error {
	if len(input) != 7 {
		return fmt.Errorf("need 7 sticks to form a number: you gave me %d", len(input))
	}

	err := l.displayPinA.Write(rpi.Value(input[0]))
	if err != nil {
		return err
	}
	err = l.displayPinB.Write(rpi.Value(input[1]))
	if err != nil {
		return err
	}
	err = l.displayPinC.Write(rpi.Value(input[2]))
	if err != nil {
		return err
	}
	err = l.displayPinD.Write(rpi.Value(input[3]))
	if err != nil {
		return err
	}
	err = l.displayPinE.Write(rpi.Value(input[4]))
	if err != nil {
		return err
	}
	err = l.displayPinF.Write(rpi.Value(input[5]))
	if err != nil {
		return err
	}
	err = l.displayPinG.Write(rpi.Value(input[6]))
	if err != nil {
		return err
	}
	return nil
}

//screen clear, shutdown all the lighted sticks
func (l *Led3461BS) Dark() error {

	err := l.displayPinA.Write(rpi.HIGH)
	if err != nil {
		return err
	}
	err = l.displayPinB.Write(rpi.HIGH)
	if err != nil {
		return err
	}
	err = l.displayPinC.Write(rpi.HIGH)
	if err != nil {
		return err
	}
	err = l.displayPinD.Write(rpi.HIGH)
	if err != nil {
		return err
	}
	err = l.displayPinE.Write(rpi.HIGH)
	if err != nil {
		return err
	}
	err = l.displayPinF.Write(rpi.HIGH)
	if err != nil {
		return err
	}
	err = l.displayPinG.Write(rpi.HIGH)
	if err != nil {
		return err
	}
	err = l.dpPin.Write(rpi.HIGH)
	if err != nil {
		return err
	}

	err = l.controlPin1.Write(rpi.LOW)
	if err != nil {
		return err
	}
	err = l.controlPin2.Write(rpi.LOW)
	if err != nil {
		return err
	}
	err = l.controlPin3.Write(rpi.LOW)
	if err != nil {
		return err
	}
	err = l.controlPin4.Write(rpi.LOW)
	if err != nil {
		return err
	}

	return nil
}

//as the method says
func (l *Led3461BS) LightUpDot(index int) error {

	err := l.controlPin1.Write(rpi.LOW)
	if err != nil {
		return err
	}
	err = l.controlPin2.Write(rpi.LOW)
	if err != nil {
		return err
	}
	err = l.controlPin3.Write(rpi.LOW)
	if err != nil {
		return err
	}
	err = l.controlPin4.Write(rpi.LOW)
	if err != nil {
		return err
	}

	switch index {
	case 1:
		err = l.controlPin1.Write(rpi.HIGH)
		if err != nil {
			return err
		}
	case 2:
		err = l.controlPin2.Write(rpi.HIGH)
		if err != nil {
			return err
		}
	case 3:
		err = l.controlPin3.Write(rpi.HIGH)
		if err != nil {
			return err
		}
	case 4:
		err = l.controlPin4.Write(rpi.HIGH)
		if err != nil {
			return err
		}
	default:
		return nil
	}
	return l.dpPin.Write(rpi.LOW)
}

//more flexible control, didn't wrote it well
//use with caution
//pos1~pos4: what should these positions show(char "H" is 10)
//dots:[true,true,true,true]=all dots light up...
// for example, "7.3H.H.", pass in these: 7,3,10,10,[4]bool{true, false, true, true}
func (l *Led3461BS) FlexibleControl(pos1, pos2, pos3, pos4 int, dots [4]bool, durationSec int) {

	timer := time.NewTimer(time.Duration(durationSec) * time.Second)
	c := make(chan int, 0)
	go func() {
		for ; ; {
			select {
			case <-c:
				return
			default:
			}
			for i := 0; i < 4; i++ {
				_ = l.Dark()
				if dots[i] {
					err := l.LightUpDot(i + 1)
					if err != nil {
						panic(err)
					}
				}
			}
			err := l.DisplaySingleChar(1, pos1)
			if err != nil {
				panic(err)
			}
			err = l.DisplaySingleChar(2, pos2)
			if err != nil {
				panic(err)
			}
			err = l.DisplaySingleChar(3, pos3)
			if err != nil {
				panic(err)
			}
			err = l.DisplaySingleChar(4, pos4)
			if err != nil {
				panic(err)
			}
		}
	}()

	//Display over
	select {
	case <-timer.C:
		c<-1
	}
	_ = l.Dark()
}
