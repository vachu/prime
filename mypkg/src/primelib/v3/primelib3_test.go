package primelib_test

import ("os"; "testing"; "primelib")

type testArgs struct {
    from, to uint32
}

var testData = []testArgs {
    //{1, 100},
    //{1, 10},
    //{90, 100},
    {5, 10},
}

func TestWritePrimesBetween(t *testing.T) {
    t.Log("From TestWritePrimesBetween...\n")
    for _, args := range testData {
        t.Logf("==== from = %d; to = %d ====\n", args.from, args.to)
        primelib.WritePrimesBetween(os.Stdout, args.from, args.to, "\n")
    }
}
