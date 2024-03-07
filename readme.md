[![Go](https://github.com/jacekdobrowolski/simple_scraper/actions/workflows/go.yml/badge.svg)](https://github.com/jacekdobrowolski/simple_scraper/actions/workflows/go.yml)

# Simple scraper

Includes concurrency and simple cache.
Made as an interview assignment.

Output:
```bash
requesting https://pl.wikipedia.org/wiki/Niemcy
requesting https://pl.wikipedia.org/wiki/Polska
requesting https://pl.wikipedia.org/wiki/Francja
done https://pl.wikipedia.org/wiki/Niemcy [{się 242} {z 315} {na 331} {i 459} {w 920}]
requesting https://pl.wikipedia.org/wiki/Polska
done https://pl.wikipedia.org/wiki/Polska [{się 250} {na 381} {z 442} {i 786} {w 1324}]
from cache https://pl.wikipedia.org/wiki/Polska
from cache https://pl.wikipedia.org/wiki/Polska
from cache https://pl.wikipedia.org/wiki/Polska
done https://pl.wikipedia.org/wiki/Francja [{z 97} {edytuj 107} {na 110} {i 159} {w 273}]
done https://pl.wikipedia.org/wiki/Polska [{się 250} {na 381} {z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Polska [{się 250} {na 381} {z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Polska [{się 250} {na 381} {z 442} {i 786} {w 1324}]
done https://pl.wikipedia.org/wiki/Polska [{się 250} {na 381} {z 442} {i 786} {w 1324}]
```