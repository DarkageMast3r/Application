CREATE TABLE IF NOT EXISTS signals (
    id SERIAL PRIMARY KEY,
    client_id UUID NOT NULL,
    type VARCHAR(50) NOT NULL,
    waarde DECIMAL(10,2) NOT NULL,
    tijdstip TIMESTAMP WITH TIME ZONE NOT NULL,
    bron VARCHAR(100) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS classifications (
    id SERIAL PRIMARY KEY,
    client_id UUID NOT NULL,
    categorie VARCHAR(100) NOT NULL,
    ernst VARCHAR(50) NOT NULL,
    motivatie TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS assessments (
    id SERIAL PRIMARY KEY,
    client_id UUID NOT NULL,
    conclusie TEXT NOT NULL,
    urgentie VARCHAR(50) NOT NULL,
    gevalideerd_door VARCHAR(100) NOT NULL,
    tijdstip TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS client_conditions (
    id SERIAL PRIMARY KEY,
    toestand_id UUID NOT NULL UNIQUE,
    client_id UUID NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'actief',
    tijdstip_registratie TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Indexes voor snelheid
CREATE INDEX IF NOT EXISTS idx_signals_client_id ON signals(client_id);
CREATE INDEX IF NOT EXISTS idx_signals_tijdstip ON signals(tijdstip);
CREATE INDEX IF NOT EXISTS idx_classifications_client_id ON classifications(client_id);
CREATE INDEX IF NOT EXISTS idx_assessments_client_id ON assessments(client_id);
CREATE INDEX IF NOT EXISTS idx_client_conditions_client_id ON client_conditions(client_id);
CREATE INDEX IF NOT EXISTS idx_client_conditions_toestand_id ON client_conditions(toestand_id);

-- Testdata
INSERT INTO signals (client_id, type, waarde, tijdstip, bron) VALUES 
    ('123e4567-e89b-12d3-a456-426614174000', 'hartslag', 72.0, NOW() - INTERVAL '1 hour', 'sensor'),
    ('123e4567-e89b-12d3-a456-426614174000', 'bloeddruk', 120.0, NOW() - INTERVAL '30 minutes', 'sensor')
ON CONFLICT DO NOTHING;

INSERT INTO classifications (client_id, categorie, ernst, motivatie) VALUES 
    ('123e4567-e89b-12d3-a456-426614174000', 'cardiovasculair', 'normaal', 'Vitalen binnen normale grenzen')
ON CONFLICT DO NOTHING;

INSERT INTO assessments (client_id, conclusie, urgentie, gevalideerd_door, tijdstip) VALUES 
    ('123e4567-e89b-12d3-a456-426614174000', 'Client in stabiele conditie', 'laag', 'Dr. van der Berg', NOW() - INTERVAL '15 minutes')
ON CONFLICT DO NOTHING;

INSERT INTO client_conditions (toestand_id, client_id, status) VALUES 
    ('123e4567-e89b-12d3-a456-426614174000', '123e4567-e89b-12d3-a456-426614174000', 'actief')
ON CONFLICT DO NOTHING;