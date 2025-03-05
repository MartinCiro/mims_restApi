DO $$ 
DECLARE 
    modelo TEXT;
    secuencia TEXT;
BEGIN
    -- Iterar sobre todas las tablas que tienen una columna 'id'
    FOR modelo IN 
        SELECT table_name 
        FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name IN ('estado', 'pedidos', 'permiso', 'rol', 'usuarios') -- Modelos de Prisma
    LOOP
        -- Obtener el nombre de la secuencia asociada a la columna 'id' de cada tabla
        SELECT pg_get_serial_sequence(modelo, 'id') INTO secuencia;

        -- Si existe una secuencia, ajustar el valor al máximo id existente
        IF secuencia IS NOT NULL THEN
            EXECUTE format('
                SELECT setval(%L, COALESCE((SELECT MAX(id) FROM %I), 1), false);
            ', secuencia, modelo);
        END IF;
    END LOOP;
END $$;
