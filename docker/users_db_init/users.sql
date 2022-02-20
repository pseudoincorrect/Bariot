create table users (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at timestamp with time zone DEFAULT now(),
  email character varying(255) NOT NULL UNIQUE,
  full_name character varying(255) NOT NULL,
  hash_pass character varying(255) NOT NULL,
  metadata json
);
