create table things (
  id uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  created_at timestamp with time zone DEFAULT now(),
  key character varying(4096) NOT NULL,
  name character varying(255) NOT NULL,
  user_id uuid NOT NULL,
  metadata json
);
