# go montecarlo

This go package comtains a variety of Monte Carlo sampling tools I use for hydrologic modelling parameter estimation.

Two basic sampling plans are provided, the Halton digital sequence (Faure and Lemieux, 2008; Lemieux, 2009) and the Latin Hypercube (Lemieux, 2009).

Also included in this pack is a set of distribution transforms that maps the uniform distribution U[0,1) to either:
* the Johnson bounded with mode m (Law, 2007)
* generalized trapezoid
* triangle

## dependencies:

* mmaths (https://github.com/maseology/mmaths)

## References

Faure, H., and C. Lemieux. 2008. Generalized Halton Sequences in 2008: A Comparative Study. 30pp.

Law, A.M., 2007. Simulation Modeling and Analysis. McGraw-Hill, fourth ed. New York. 768pp.

Lemieux, C., Monte Carlo and Quasi-Monte Carlo Sampling. Springer Science. 2009. 373pp.