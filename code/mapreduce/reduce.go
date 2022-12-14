package mapreduce

import (
	"fmt"
	"log"
	"math/rand"
)

// Reduce : per ogni cluster calcola la media del nuovo centroide
func (a *API) Reduce(clusters []Clusters, centroid *[]Centroids) error {
	fmt.Println(" ")
	log.Printf(" ** Reduce phase **")
	log.Printf("Numero di cluster assegnati %d", len(clusters))

	for ii := range clusters {
		log.Printf("Cluster %d con %d punti", clusters[ii].Centroid.Index, len(clusters[ii].PointsData))
	}

	//Crea centroid
	centroids := make([]Centroids, len(clusters))
	var lenPoint int

	if len(clusters) != 0 {
		for ii := range clusters {
			if len(clusters[ii].PointsData) != 0 {

				lenPoint = len(clusters[ii].PointsData[0].Point)

				p := make([][]float64, lenPoint)

				// Calcola la media di ogni cluster
				for j := range clusters[ii].PointsData {
					for k := range clusters[ii].PointsData[j].Point {
						p[k] = append(p[k], clusters[ii].PointsData[j].Point[k])
					}
				}

				var mean []float64
				for k := range p {
					var sum float64
					for j := range p[k] {
						sum += p[k][j]
					}
					var op = sum / float64(len(p[k]))
					mean = append(mean, op)
				}

				centroids[ii].Index = ii
				centroids[ii].Centroid = mean
			} else {
				centroids[ii].Index = ii
				centroids[ii].Centroid = clusters[ii].Centroid.Centroid
			}
		}
	}
	*centroid = centroids

	return nil
}

// ReduceKMeans
// We reduce the first element of two value pairs
// by choosing one of these elements with probability proportional to the second
// element in each pair, and reduce the second element of the pairs by summation.
func (a *API) ReduceKMeans(points []Points, centroid *[]Centroids) error {
	fmt.Println(" ")
	log.Printf(" ** Reduce phase K-means++ **")
	log.Printf("Numero di punti assegnati %d", len(points))

	//Crea centroid
	var centroids []Centroids

	// Calcola il nuovo centroide in base ad una distanza casuale calcolata come somma delle distanze dei punti

	distance := make([]float64, len(points))
	var c Centroids
	var sum float64
	jj := 0

	for pp := range points {
		distance[pp] = points[pp].Distance * points[pp].Distance
		sum += distance[pp]
	}

	//Trova una distanza casuale moltiplicando un valore random con la somma delle distanze
	randomDistance := rand.Float64() * sum

	// Assegna i centroidi.
	for sum = distance[0]; sum < randomDistance; sum += distance[jj] {
		jj++
	}

	c.Index = len(points[0].Centroids) + 1
	c.Centroid = points[jj].Point
	centroids = append(centroids, c)

	*centroid = centroids
	return nil
}
