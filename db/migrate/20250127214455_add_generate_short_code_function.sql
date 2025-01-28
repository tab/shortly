-- +goose Up
CREATE FUNCTION generate_short_code()
RETURNS text AS $$
DECLARE
  val bigint;
  code text;
BEGIN
SELECT nextval('short_code_seq') INTO val;
  code := base62(val);
  IF length(code) > 8 THEN code := substring(code for 8);
END IF;
  WHILE length(code) < 8 LOOP
    code := '0' || code;
  END LOOP;
RETURN code;
END;
$$ LANGUAGE plpgsql;

-- +goose Down
DROP FUNCTION generate_short_code();
