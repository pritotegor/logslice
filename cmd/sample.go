package cmd

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"

	"github.com/spf13/cobra"
)

var sampleCmd = &cobra.Command{
	Use:   "sample",
	Short: "Sample a fraction or fixed count of log lines",
	RunE:  RunSample,
}

func init() {
	sampleCmd.Flags().Float64P("rate", "r", 0.0, "Sampling rate between 0.0 and 1.0 (e.g. 0.1 for 10%)")
	sampleCmd.Flags().IntP("count", "n", 0, "Sample exactly N random lines (reservoir sampling)")
	sampleCmd.Flags().StringP("input", "i", "", "Input file (default: stdin)")
	sampleCmd.Flags().Int64("seed", 0, "Random seed for reproducibility (0 = random)")
	rootCmd.AddCommand(sampleCmd)
}

func RunSample(cmd *cobra.Command, args []string) error {
	rate, _ := cmd.Flags().GetFloat64("rate")
	count, _ := cmd.Flags().GetInt("count")
	inputFile, _ := cmd.Flags().GetString("input")
	seed, _ := cmd.Flags().GetInt64("seed")

	if rate == 0.0 && count == 0 {
		return fmt.Errorf("either --rate or --count must be specified")
	}
	if rate != 0.0 && (rate < 0.0 || rate > 1.0) {
		return fmt.Errorf("--rate must be between 0.0 and 1.0")
	}

	r := newRand(seed)

	in, err := openInput(inputFile)
	if err != nil {
		return err
	}
	if f, ok := in.(*os.File); ok && f != os.Stdin {
		defer f.Close()
	}

	scanner := bufio.NewScanner(in)

	if count > 0 {
		return reservoirSample(scanner, count, r)
	}
	return rateSample(scanner, rate, r)
}

func rateSample(scanner *bufio.Scanner, rate float64, r *rand.Rand) error {
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if r.Float64() < rate {
			fmt.Println(line)
		}
	}
	return scanner.Err()
}

func reservoirSample(scanner *bufio.Scanner, n int, r *rand.Rand) error {
	reservoir := make([]string, 0, n)
	i := 0
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			continue
		}
		if len(reservoir) < n {
			reservoir = append(reservoir, line)
		} else {
			j := r.Intn(i + 1)
			if j < n {
				reservoir[j] = line
			}
		}
		i++
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	for _, line := range reservoir {
		fmt.Println(line)
	}
	return nil
}

func newRand(seed int64) *rand.Rand {
	if seed == 0 {
		seed = rand.Int63()
	}
	return rand.New(rand.NewSource(seed))
}
