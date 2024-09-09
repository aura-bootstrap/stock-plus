package function

type Function interface {
	String() string
	Main(input <-chan string, output chan<- string)
}
