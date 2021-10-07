ALTER Table
  movies
ADD
  CONSTRAINT movies_duration_check CHECK (duration >= 0);
ALTER Table
  movies
ADD
  CONSTRAINT movies_year_check CHECK (
    year BETWEEN 1888
    AND date_part('year', now())
  );
ALTER Table
  movies
ADD
  CONSTRAINT genres_length_check CHECK (
    array_length(genres, 1) BETWEEN 1
    and 5
  );