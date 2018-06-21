# go montecarlo

This go package comtains a variety of Monte Carlo sampling tools I use for hydrologic modelling parameter estimation.

Two basic sampling plans are provided, the Halton digital sequence (Faure and Lemieux, 2008; Lemieux, 2009) and the Latin Hypercube (Lemieux, 2009).

# Also included: 

A set of distribution transforms that maps the uniform distribution U[0,1) to either:
* the Johnson bounded with mode m (Law, 2007)
* generalized trapezoid
* triangle

A set of joint distribution transforms:
* Nested distributions (i.e., 0.0 <= u1 <= u2 <= 1.0)
* Three symmetric and invertable copulae: elliptical, Franks archimedean, diagonal band (from Kurowicka and Cooke, 2006)

## dependencies:

* mmaths (https://github.com/maseology/mmaths)

## References

Faure, H., and C. Lemieux, 2008. Generalized Halton Sequences in 2008: A Comparative Study. 30pp.

Kurowicka, D. and R. Cooke, 2006. Uncertainty Analysis with High Dimensional Dependence Modelling. John Wiley & Sons, Ltd. 284pp.

Law, A.M., 2007. Simulation Modeling and Analysis. McGraw-Hill, fourth ed. New York. 768pp.

Lemieux, C., 2009. Monte Carlo and Quasi-Monte Carlo Sampling. Springer Science. 373pp.