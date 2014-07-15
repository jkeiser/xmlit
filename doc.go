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

  decoder := xmlit.DecodeElements(reader, Customer{}, Order{})
  for decoder.HasNext() {
    fmt.Println(decoder.Next())
  }
  if decoder.Error() {
    // process error
  }
*/
package xmlit
