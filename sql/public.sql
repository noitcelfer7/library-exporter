CREATE TABLE public.authors (
    uuid         uuid         DEFAULT gen_random_uuid() NOT NULL,
    first_name   varchar(32)  							NOT NULL,
    last_name    varchar(32)  							NOT NULL,

    CONSTRAINT authors_pkey PRIMARY KEY (uuid),

    CONSTRAINT authors_ukey UNIQUE (first_name, last_name)
);

CREATE TABLE public.readers (
    uuid          uuid         DEFAULT gen_random_uuid() NOT NULL,
    first_name    varchar(32)  							 NOT NULL,
    last_name     varchar(32)  							 NOT NULL,
    phone_number  varchar(32)  							 NOT NULL,

    CONSTRAINT readers_pkey PRIMARY KEY (uuid),

    CONSTRAINT readers_ukey1 UNIQUE (first_name, last_name),
    CONSTRAINT readers_ukey2 UNIQUE (phone_number)
);

CREATE TABLE public.books (
    uuid   uuid         DEFAULT gen_random_uuid() NOT NULL,
    isbn   varchar(32)  						  NOT NULL,
    title  varchar(32)  						  NOT NULL,

    CONSTRAINT books_pkey PRIMARY KEY (uuid),

    CONSTRAINT books_ukey1 UNIQUE (isbn),
    CONSTRAINT books_ukey2 UNIQUE (title)
);

CREATE TABLE public.genres (
    uuid   uuid         DEFAULT gen_random_uuid() NOT NULL,
    title  varchar(32)  						  NOT NULL,

    CONSTRAINT genres_pkey PRIMARY KEY (uuid),

    CONSTRAINT genres_ukey UNIQUE (title)
);

CREATE TABLE public.issues (
    uuid         uuid         DEFAULT gen_random_uuid() NOT NULL,
    book_uuid    uuid,
    issue_date   date         DEFAULT now() NOT NULL,
    period       int          				NOT NULL,
    reader_uuid  uuid,
    return_date  date,

    CONSTRAINT issues_pkey PRIMARY KEY (uuid),

    CONSTRAINT issues_book_fk   FOREIGN KEY (book_uuid)   REFERENCES public.books(uuid),
    CONSTRAINT issues_reader_fk FOREIGN KEY (reader_uuid) REFERENCES public.readers(uuid),

    CONSTRAINT genres_ukey1 UNIQUE (book_uuid, issue_date, reader_uuid)
);

CREATE TABLE public.authors_to_books (
    author_uuid  uuid NOT NULL,
    book_uuid    uuid NOT NULL,

    CONSTRAINT authors_to_books_pk PRIMARY KEY (author_uuid, book_uuid),

    CONSTRAINT authors_to_books_author_fk FOREIGN KEY (author_uuid) REFERENCES public.authors(uuid),
    CONSTRAINT authors_to_books_book_fk   FOREIGN KEY (book_uuid)   REFERENCES public.books(uuid)
);

CREATE TABLE public.books_to_genres (
    book_uuid   uuid NOT NULL,
    genre_uuid  uuid NOT NULL,

    CONSTRAINT books_to_genres_pk PRIMARY KEY (book_uuid, genre_uuid),

    CONSTRAINT books_to_genres_book_fk  FOREIGN KEY (book_uuid)  REFERENCES public.books(uuid),
    CONSTRAINT books_to_genres_genre_fk FOREIGN KEY (genre_uuid) REFERENCES public.genres(uuid)
);
