/*
Allows the efficient decoding of large XML streams.

Iterates through the XML file, finding and parsing matching elements from the
stream and returning them, in order. Uses xml.Decode, so you can use structs
for meaningful and immediately usable results.

  struct Customer {
    Id      int
    Name    string
    Address string
  }
  struct Order {
    CustomerId int
    Total float64
  }

  iterator := xmlit.DecodeElements(reader, Customer{}, Order{})
  iterator.Each(func (item interface{}) { fmt.Println(item) })
*/
package xmlit
