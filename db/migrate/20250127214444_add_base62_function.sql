-- +goose Up
CREATE FUNCTION base62(val bigint)
  RETURNS text AS $func$
DECLARE
alpha text := '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz';
  code text := '';
  remainder int;
  tmp bigint := val;
BEGIN
  IF tmp = 0 THEN
    RETURN '0';
END IF;

  WHILE tmp > 0 LOOP
    remainder := tmp % 62;
    code := substr(alpha, remainder + 1, 1) || code;
    tmp := tmp / 62;
END LOOP;

RETURN code;
END;
$func$ LANGUAGE plpgsql IMMUTABLE STRICT;

-- +goose Down
DROP FUNCTION base62(bigint);
