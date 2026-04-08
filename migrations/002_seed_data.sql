-- Seed data for RateYourProduction
-- Real theatre works, people, companies, venues, and productions

BEGIN;

-- ============================================================
-- GENRES
-- ============================================================
INSERT INTO genres (name, slug) VALUES
    ('Tragedy', 'tragedy'),
    ('Comedy', 'comedy'),
    ('Drama', 'drama'),
    ('Musical Theatre', 'musical-theatre'),
    ('Opera', 'opera'),
    ('Absurdist', 'absurdist'),
    ('Historical', 'historical'),
    ('Romance', 'romance'),
    ('Farce', 'farce'),
    ('Thriller', 'thriller'),
    ('Experimental', 'experimental'),
    ('Tragicomedy', 'tragicomedy')
ON CONFLICT DO NOTHING;

-- ============================================================
-- PEOPLE
-- ============================================================
INSERT INTO people (name, slug) VALUES
    -- Playwrights
    ('William Shakespeare', 'william-shakespeare'),
    ('Arthur Miller', 'arthur-miller'),
    ('Tennessee Williams', 'tennessee-williams'),
    ('Samuel Beckett', 'samuel-beckett'),
    ('Anton Chekhov', 'anton-chekhov'),
    ('Henrik Ibsen', 'henrik-ibsen'),
    ('Tom Stoppard', 'tom-stoppard'),
    ('Caryl Churchill', 'caryl-churchill'),
    ('August Wilson', 'august-wilson'),
    ('Tony Kushner', 'tony-kushner'),
    -- Composers / Lyricists
    ('Stephen Sondheim', 'stephen-sondheim'),
    ('Andrew Lloyd Webber', 'andrew-lloyd-webber'),
    ('Lin-Manuel Miranda', 'lin-manuel-miranda'),
    ('Giuseppe Verdi', 'giuseppe-verdi'),
    ('Giacomo Puccini', 'giacomo-puccini'),
    ('James Lapine', 'james-lapine'),
    ('Tim Rice', 'tim-rice'),
    -- Directors / Actors
    ('Sam Mendes', 'sam-mendes'),
    ('Marianne Elliott', 'marianne-elliott'),
    ('Ivo van Hove', 'ivo-van-hove'),
    ('Robert Icke', 'robert-icke'),
    ('Ian McKellen', 'ian-mckellen'),
    ('Cate Blanchett', 'cate-blanchett'),
    ('Denzel Washington', 'denzel-washington'),
    ('Audra McDonald', 'audra-mcdonald')
ON CONFLICT DO NOTHING;

-- ============================================================
-- WORKS
-- ============================================================
INSERT INTO works (slug, title, normalized_title, type, description, premiere_year) VALUES
    -- Plays
    ('hamlet', 'Hamlet', 'hamlet', 'play',
     'The Prince of Denmark seeks revenge after his father is murdered by his uncle, who then marries his mother and claims the throne.', 1601),
    ('death-of-a-salesman', 'Death of a Salesman', 'death of a salesman', 'play',
     'Willy Loman, an aging travelling salesman, confronts the gap between his dreams and reality as his life unravels.', 1949),
    ('a-streetcar-named-desire', 'A Streetcar Named Desire', 'a streetcar named desire', 'play',
     'Blanche DuBois arrives at her sister Stella''s home in New Orleans, where she clashes with Stella''s husband Stanley Kowalski.', 1947),
    ('waiting-for-godot', 'Waiting for Godot', 'waiting for godot', 'play',
     'Two men wait endlessly and in vain for someone named Godot in a tragicomic exploration of existence and meaning.', 1953),
    ('the-cherry-orchard', 'The Cherry Orchard', 'the cherry orchard', 'play',
     'An aristocratic Russian family faces the loss of their beloved estate and its cherry orchard to pay off debts.', 1904),
    ('a-dolls-house', 'A Doll''s House', 'a dolls house', 'play',
     'Nora Helmer discovers that her life of domestic comfort has been built on deception and that she must find her own identity.', 1879),
    ('rosencrantz-and-guildenstern-are-dead', 'Rosencrantz and Guildenstern Are Dead', 'rosencrantz and guildenstern are dead', 'play',
     'Two minor characters from Hamlet find themselves caught up in events beyond their understanding in this existential comedy.', 1966),
    ('top-girls', 'Top Girls', 'top girls', 'play',
     'Marlene celebrates her promotion at a dinner party with famous women from history, then confronts the cost of her success.', 1982),
    ('fences', 'Fences', 'fences', 'play',
     'Troy Maxson, a former Negro Leagues baseball player, struggles with his role as a father and husband in 1950s Pittsburgh.', 1985),
    ('angels-in-america-millennium-approaches', 'Angels in America: Millennium Approaches', 'angels in america millennium approaches', 'play',
     'Interconnected stories of New Yorkers grappling with AIDS, politics, and identity in Reagan-era America.', 1991),
    ('king-lear', 'King Lear', 'king lear', 'play',
     'An aging king divides his realm among his daughters based on their flattery, leading to betrayal, madness, and devastation.', 1606),
    ('the-crucible', 'The Crucible', 'the crucible', 'play',
     'A dramatization of the Salem witch trials used as an allegory for McCarthyism in 1950s America.', 1953),
    -- Musicals
    ('hamilton', 'Hamilton', 'hamilton', 'musical',
     'The life of American Founding Father Alexander Hamilton told through hip-hop, R&B, and traditional show tunes.', 2015),
    ('sweeney-todd', 'Sweeney Todd: The Demon Barber of Fleet Street', 'sweeney todd the demon barber of fleet street', 'musical',
     'A wrongfully exiled barber returns to London seeking revenge against the judge who destroyed his family.', 1979),
    ('into-the-woods', 'Into the Woods', 'into the woods', 'musical',
     'Several fairy-tale characters embark on intertwining quests to fulfill their wishes, only to face the consequences.', 1987),
    ('the-phantom-of-the-opera', 'The Phantom of the Opera', 'the phantom of the opera', 'musical',
     'A disfigured musical genius lurks beneath the Paris Opera House, obsessed with a young soprano.', 1986),
    -- Operas
    ('la-traviata', 'La Traviata', 'la traviata', 'opera',
     'A Parisian courtesan sacrifices her own happiness for the sake of her lover''s family honour.', 1853),
    ('la-boheme', 'La Bohème', 'la boheme', 'opera',
     'Impoverished young artists and musicians struggle to survive in 1830s Paris as love blossoms and fades.', 1896)
ON CONFLICT DO NOTHING;

-- ============================================================
-- WORK CREATORS
-- ============================================================
INSERT INTO work_creators (work_id, person_id, role_type) VALUES
    -- Shakespeare plays
    ((SELECT id FROM works WHERE slug = 'hamlet'), (SELECT id FROM people WHERE slug = 'william-shakespeare'), 'playwright'),
    ((SELECT id FROM works WHERE slug = 'king-lear'), (SELECT id FROM people WHERE slug = 'william-shakespeare'), 'playwright'),
    -- Arthur Miller
    ((SELECT id FROM works WHERE slug = 'death-of-a-salesman'), (SELECT id FROM people WHERE slug = 'arthur-miller'), 'playwright'),
    ((SELECT id FROM works WHERE slug = 'the-crucible'), (SELECT id FROM people WHERE slug = 'arthur-miller'), 'playwright'),
    -- Tennessee Williams
    ((SELECT id FROM works WHERE slug = 'a-streetcar-named-desire'), (SELECT id FROM people WHERE slug = 'tennessee-williams'), 'playwright'),
    -- Samuel Beckett
    ((SELECT id FROM works WHERE slug = 'waiting-for-godot'), (SELECT id FROM people WHERE slug = 'samuel-beckett'), 'playwright'),
    -- Chekhov
    ((SELECT id FROM works WHERE slug = 'the-cherry-orchard'), (SELECT id FROM people WHERE slug = 'anton-chekhov'), 'playwright'),
    -- Ibsen
    ((SELECT id FROM works WHERE slug = 'a-dolls-house'), (SELECT id FROM people WHERE slug = 'henrik-ibsen'), 'playwright'),
    -- Stoppard
    ((SELECT id FROM works WHERE slug = 'rosencrantz-and-guildenstern-are-dead'), (SELECT id FROM people WHERE slug = 'tom-stoppard'), 'playwright'),
    -- Churchill
    ((SELECT id FROM works WHERE slug = 'top-girls'), (SELECT id FROM people WHERE slug = 'caryl-churchill'), 'playwright'),
    -- August Wilson
    ((SELECT id FROM works WHERE slug = 'fences'), (SELECT id FROM people WHERE slug = 'august-wilson'), 'playwright'),
    -- Kushner
    ((SELECT id FROM works WHERE slug = 'angels-in-america-millennium-approaches'), (SELECT id FROM people WHERE slug = 'tony-kushner'), 'playwright'),
    -- Hamilton
    ((SELECT id FROM works WHERE slug = 'hamilton'), (SELECT id FROM people WHERE slug = 'lin-manuel-miranda'), 'composer'),
    ((SELECT id FROM works WHERE slug = 'hamilton'), (SELECT id FROM people WHERE slug = 'lin-manuel-miranda'), 'lyricist'),
    ((SELECT id FROM works WHERE slug = 'hamilton'), (SELECT id FROM people WHERE slug = 'lin-manuel-miranda'), 'book_writer'),
    -- Sweeney Todd
    ((SELECT id FROM works WHERE slug = 'sweeney-todd'), (SELECT id FROM people WHERE slug = 'stephen-sondheim'), 'composer'),
    ((SELECT id FROM works WHERE slug = 'sweeney-todd'), (SELECT id FROM people WHERE slug = 'stephen-sondheim'), 'lyricist'),
    -- Into the Woods
    ((SELECT id FROM works WHERE slug = 'into-the-woods'), (SELECT id FROM people WHERE slug = 'stephen-sondheim'), 'composer'),
    ((SELECT id FROM works WHERE slug = 'into-the-woods'), (SELECT id FROM people WHERE slug = 'stephen-sondheim'), 'lyricist'),
    ((SELECT id FROM works WHERE slug = 'into-the-woods'), (SELECT id FROM people WHERE slug = 'james-lapine'), 'book_writer'),
    -- Phantom
    ((SELECT id FROM works WHERE slug = 'the-phantom-of-the-opera'), (SELECT id FROM people WHERE slug = 'andrew-lloyd-webber'), 'composer'),
    ((SELECT id FROM works WHERE slug = 'the-phantom-of-the-opera'), (SELECT id FROM people WHERE slug = 'tim-rice'), 'lyricist'),
    -- La Traviata
    ((SELECT id FROM works WHERE slug = 'la-traviata'), (SELECT id FROM people WHERE slug = 'giuseppe-verdi'), 'composer'),
    -- La Bohème
    ((SELECT id FROM works WHERE slug = 'la-boheme'), (SELECT id FROM people WHERE slug = 'giacomo-puccini'), 'composer')
ON CONFLICT DO NOTHING;

-- ============================================================
-- WORK GENRES
-- ============================================================
INSERT INTO work_genres (work_id, genre_id) VALUES
    -- Hamlet
    ((SELECT id FROM works WHERE slug = 'hamlet'), (SELECT id FROM genres WHERE slug = 'tragedy')),
    ((SELECT id FROM works WHERE slug = 'hamlet'), (SELECT id FROM genres WHERE slug = 'drama')),
    -- Death of a Salesman
    ((SELECT id FROM works WHERE slug = 'death-of-a-salesman'), (SELECT id FROM genres WHERE slug = 'tragedy')),
    ((SELECT id FROM works WHERE slug = 'death-of-a-salesman'), (SELECT id FROM genres WHERE slug = 'drama')),
    -- Streetcar
    ((SELECT id FROM works WHERE slug = 'a-streetcar-named-desire'), (SELECT id FROM genres WHERE slug = 'drama')),
    ((SELECT id FROM works WHERE slug = 'a-streetcar-named-desire'), (SELECT id FROM genres WHERE slug = 'tragedy')),
    -- Waiting for Godot
    ((SELECT id FROM works WHERE slug = 'waiting-for-godot'), (SELECT id FROM genres WHERE slug = 'absurdist')),
    ((SELECT id FROM works WHERE slug = 'waiting-for-godot'), (SELECT id FROM genres WHERE slug = 'tragicomedy')),
    -- Cherry Orchard
    ((SELECT id FROM works WHERE slug = 'the-cherry-orchard'), (SELECT id FROM genres WHERE slug = 'drama')),
    ((SELECT id FROM works WHERE slug = 'the-cherry-orchard'), (SELECT id FROM genres WHERE slug = 'comedy')),
    -- A Doll's House
    ((SELECT id FROM works WHERE slug = 'a-dolls-house'), (SELECT id FROM genres WHERE slug = 'drama')),
    -- Rosencrantz
    ((SELECT id FROM works WHERE slug = 'rosencrantz-and-guildenstern-are-dead'), (SELECT id FROM genres WHERE slug = 'absurdist')),
    ((SELECT id FROM works WHERE slug = 'rosencrantz-and-guildenstern-are-dead'), (SELECT id FROM genres WHERE slug = 'comedy')),
    -- Top Girls
    ((SELECT id FROM works WHERE slug = 'top-girls'), (SELECT id FROM genres WHERE slug = 'drama')),
    ((SELECT id FROM works WHERE slug = 'top-girls'), (SELECT id FROM genres WHERE slug = 'experimental')),
    -- Fences
    ((SELECT id FROM works WHERE slug = 'fences'), (SELECT id FROM genres WHERE slug = 'drama')),
    ((SELECT id FROM works WHERE slug = 'fences'), (SELECT id FROM genres WHERE slug = 'tragedy')),
    -- Angels in America
    ((SELECT id FROM works WHERE slug = 'angels-in-america-millennium-approaches'), (SELECT id FROM genres WHERE slug = 'drama')),
    ((SELECT id FROM works WHERE slug = 'angels-in-america-millennium-approaches'), (SELECT id FROM genres WHERE slug = 'historical')),
    -- King Lear
    ((SELECT id FROM works WHERE slug = 'king-lear'), (SELECT id FROM genres WHERE slug = 'tragedy')),
    ((SELECT id FROM works WHERE slug = 'king-lear'), (SELECT id FROM genres WHERE slug = 'drama')),
    -- The Crucible
    ((SELECT id FROM works WHERE slug = 'the-crucible'), (SELECT id FROM genres WHERE slug = 'drama')),
    ((SELECT id FROM works WHERE slug = 'the-crucible'), (SELECT id FROM genres WHERE slug = 'historical')),
    -- Hamilton
    ((SELECT id FROM works WHERE slug = 'hamilton'), (SELECT id FROM genres WHERE slug = 'musical-theatre')),
    ((SELECT id FROM works WHERE slug = 'hamilton'), (SELECT id FROM genres WHERE slug = 'historical')),
    -- Sweeney Todd
    ((SELECT id FROM works WHERE slug = 'sweeney-todd'), (SELECT id FROM genres WHERE slug = 'musical-theatre')),
    ((SELECT id FROM works WHERE slug = 'sweeney-todd'), (SELECT id FROM genres WHERE slug = 'thriller')),
    -- Into the Woods
    ((SELECT id FROM works WHERE slug = 'into-the-woods'), (SELECT id FROM genres WHERE slug = 'musical-theatre')),
    ((SELECT id FROM works WHERE slug = 'into-the-woods'), (SELECT id FROM genres WHERE slug = 'comedy')),
    -- Phantom
    ((SELECT id FROM works WHERE slug = 'the-phantom-of-the-opera'), (SELECT id FROM genres WHERE slug = 'musical-theatre')),
    ((SELECT id FROM works WHERE slug = 'the-phantom-of-the-opera'), (SELECT id FROM genres WHERE slug = 'romance')),
    -- La Traviata
    ((SELECT id FROM works WHERE slug = 'la-traviata'), (SELECT id FROM genres WHERE slug = 'opera')),
    ((SELECT id FROM works WHERE slug = 'la-traviata'), (SELECT id FROM genres WHERE slug = 'tragedy')),
    ((SELECT id FROM works WHERE slug = 'la-traviata'), (SELECT id FROM genres WHERE slug = 'romance')),
    -- La Bohème
    ((SELECT id FROM works WHERE slug = 'la-boheme'), (SELECT id FROM genres WHERE slug = 'opera')),
    ((SELECT id FROM works WHERE slug = 'la-boheme'), (SELECT id FROM genres WHERE slug = 'romance')),
    ((SELECT id FROM works WHERE slug = 'la-boheme'), (SELECT id FROM genres WHERE slug = 'tragedy'))
ON CONFLICT DO NOTHING;

-- ============================================================
-- COMPANIES
-- ============================================================
INSERT INTO companies (name, slug, city, country) VALUES
    ('Royal Shakespeare Company', 'royal-shakespeare-company', 'Stratford-upon-Avon', 'United Kingdom'),
    ('National Theatre', 'national-theatre', 'London', 'United Kingdom'),
    ('Stratford Festival', 'stratford-festival', 'Stratford', 'Canada'),
    ('Donmar Warehouse', 'donmar-warehouse', 'London', 'United Kingdom'),
    ('Lincoln Center Theater', 'lincoln-center-theater', 'New York', 'United States'),
    ('Roundabout Theatre Company', 'roundabout-theatre-company', 'New York', 'United States'),
    ('Almeida Theatre Company', 'almeida-theatre', 'London', 'United Kingdom'),
    ('Sydney Theatre Company', 'sydney-theatre-company', 'Sydney', 'Australia'),
    ('The Metropolitan Opera', 'metropolitan-opera', 'New York', 'United States'),
    ('Royal Opera House', 'royal-opera-house', 'London', 'United Kingdom'),
    ('Schaubühne Berlin', 'schaubuehne-berlin', 'Berlin', 'Germany'),
    ('Steppenwolf Theatre Company', 'steppenwolf-theatre', 'Chicago', 'United States')
ON CONFLICT DO NOTHING;

-- ============================================================
-- VENUES
-- ============================================================
INSERT INTO venues (name, slug, city, country) VALUES
    ('Royal Shakespeare Theatre', 'royal-shakespeare-theatre', 'Stratford-upon-Avon', 'United Kingdom'),
    ('Olivier Theatre', 'olivier-theatre', 'London', 'United Kingdom'),
    ('Lyttelton Theatre', 'lyttelton-theatre', 'London', 'United Kingdom'),
    ('Richard Rodgers Theatre', 'richard-rodgers-theatre', 'New York', 'United States'),
    ('Her Majesty''s Theatre', 'her-majestys-theatre', 'London', 'United Kingdom'),
    ('Donmar Warehouse', 'donmar-warehouse-venue', 'London', 'United Kingdom'),
    ('Almeida Theatre', 'almeida-theatre-venue', 'London', 'United Kingdom'),
    ('Festival Theatre', 'festival-theatre-stratford', 'Stratford', 'Canada'),
    ('Metropolitan Opera House', 'metropolitan-opera-house', 'New York', 'United States'),
    ('Royal Opera House', 'royal-opera-house-venue', 'London', 'United Kingdom'),
    ('Roslyn Packer Theatre', 'roslyn-packer-theatre', 'Sydney', 'Australia'),
    ('American Airlines Theatre', 'american-airlines-theatre', 'New York', 'United States')
ON CONFLICT DO NOTHING;

-- ============================================================
-- PRODUCTIONS
-- ============================================================
INSERT INTO productions (work_id, slug, company_id, venue_id, city, country, start_date, end_date, production_label) VALUES
    -- Hamlet - RSC 2016
    ((SELECT id FROM works WHERE slug = 'hamlet'),
     'hamlet-rsc-2016',
     (SELECT id FROM companies WHERE slug = 'royal-shakespeare-company'),
     (SELECT id FROM venues WHERE slug = 'royal-shakespeare-theatre'),
     'Stratford-upon-Avon', 'United Kingdom', '2016-03-17', '2016-08-13',
     'RSC 2016 Production'),

    -- Hamlet - Almeida / Robert Icke 2017
    ((SELECT id FROM works WHERE slug = 'hamlet'),
     'hamlet-almeida-2017',
     (SELECT id FROM companies WHERE slug = 'almeida-theatre'),
     (SELECT id FROM venues WHERE slug = 'almeida-theatre-venue'),
     'London', 'United Kingdom', '2017-02-17', '2017-04-08',
     'Almeida 2017 — Robert Icke'),

    -- King Lear - NT / Ian McKellen 2018
    ((SELECT id FROM works WHERE slug = 'king-lear'),
     'king-lear-nt-2018',
     (SELECT id FROM companies WHERE slug = 'national-theatre'),
     (SELECT id FROM venues WHERE slug = 'olivier-theatre'),
     'London', 'United Kingdom', '2018-07-26', '2018-11-03',
     'National Theatre 2018'),

    -- Death of a Salesman - Donmar 2019
    ((SELECT id FROM works WHERE slug = 'death-of-a-salesman'),
     'death-of-a-salesman-donmar-2019',
     (SELECT id FROM companies WHERE slug = 'donmar-warehouse'),
     (SELECT id FROM venues WHERE slug = 'donmar-warehouse-venue'),
     'London', 'United Kingdom', '2019-05-04', '2019-07-13',
     'Donmar Warehouse 2019'),

    -- A Streetcar Named Desire - NT 2014
    ((SELECT id FROM works WHERE slug = 'a-streetcar-named-desire'),
     'streetcar-nt-2014',
     (SELECT id FROM companies WHERE slug = 'national-theatre'),
     (SELECT id FROM venues WHERE slug = 'olivier-theatre'),
     'London', 'United Kingdom', '2014-07-17', '2014-11-01',
     'National Theatre 2014'),

    -- Waiting for Godot - Stratford Festival 2022
    ((SELECT id FROM works WHERE slug = 'waiting-for-godot'),
     'waiting-for-godot-stratford-2022',
     (SELECT id FROM companies WHERE slug = 'stratford-festival'),
     (SELECT id FROM venues WHERE slug = 'festival-theatre-stratford'),
     'Stratford', 'Canada', '2022-06-02', '2022-10-29',
     'Stratford Festival 2022'),

    -- The Cherry Orchard - Donmar 2023
    ((SELECT id FROM works WHERE slug = 'the-cherry-orchard'),
     'cherry-orchard-donmar-2023',
     (SELECT id FROM companies WHERE slug = 'donmar-warehouse'),
     (SELECT id FROM venues WHERE slug = 'donmar-warehouse-venue'),
     'London', 'United Kingdom', '2023-02-10', '2023-04-01',
     'Donmar Warehouse 2023'),

    -- A Doll's House - Almeida / Robert Icke 2017
    ((SELECT id FROM works WHERE slug = 'a-dolls-house'),
     'dolls-house-almeida-2017',
     (SELECT id FROM companies WHERE slug = 'almeida-theatre'),
     (SELECT id FROM venues WHERE slug = 'almeida-theatre-venue'),
     'London', 'United Kingdom', '2017-06-14', '2017-07-22',
     'Almeida 2017 — Adapted by Robert Icke'),

    -- Rosencrantz and Guildenstern Are Dead - NT 2017
    ((SELECT id FROM works WHERE slug = 'rosencrantz-and-guildenstern-are-dead'),
     'rosencrantz-nt-2017',
     (SELECT id FROM companies WHERE slug = 'national-theatre'),
     (SELECT id FROM venues WHERE slug = 'olivier-theatre'),
     'London', 'United Kingdom', '2017-02-25', '2017-05-06',
     'National Theatre 50th Anniversary Revival'),

    -- Top Girls - NT 2019
    ((SELECT id FROM works WHERE slug = 'top-girls'),
     'top-girls-nt-2019',
     (SELECT id FROM companies WHERE slug = 'national-theatre'),
     (SELECT id FROM venues WHERE slug = 'lyttelton-theatre'),
     'London', 'United Kingdom', '2019-04-04', '2019-06-22',
     'National Theatre 2019'),

    -- Fences - Roundabout 2010
    ((SELECT id FROM works WHERE slug = 'fences'),
     'fences-roundabout-2010',
     (SELECT id FROM companies WHERE slug = 'roundabout-theatre-company'),
     (SELECT id FROM venues WHERE slug = 'american-airlines-theatre'),
     'New York', 'United States', '2010-04-26', '2010-07-11',
     'Broadway Revival 2010'),

    -- Angels in America - NT 2017
    ((SELECT id FROM works WHERE slug = 'angels-in-america-millennium-approaches'),
     'angels-in-america-nt-2017',
     (SELECT id FROM companies WHERE slug = 'national-theatre'),
     (SELECT id FROM venues WHERE slug = 'lyttelton-theatre'),
     'London', 'United Kingdom', '2017-04-20', '2017-08-19',
     'National Theatre 2017 — Marianne Elliott'),

    -- The Crucible - NT 2014
    ((SELECT id FROM works WHERE slug = 'the-crucible'),
     'crucible-nt-2014',
     (SELECT id FROM companies WHERE slug = 'national-theatre'),
     (SELECT id FROM venues WHERE slug = 'olivier-theatre'),
     'London', 'United Kingdom', '2014-06-21', '2014-09-13',
     'National Theatre 2014 — Ivo van Hove'),

    -- Hamilton - Original Broadway 2015
    ((SELECT id FROM works WHERE slug = 'hamilton'),
     'hamilton-broadway-2015',
     (SELECT id FROM companies WHERE slug = 'lincoln-center-theater'),
     (SELECT id FROM venues WHERE slug = 'richard-rodgers-theatre'),
     'New York', 'United States', '2015-08-06', NULL,
     'Original Broadway Production'),

    -- Hamilton - West End 2017
    ((SELECT id FROM works WHERE slug = 'hamilton'),
     'hamilton-west-end-2017',
     NULL,
     NULL,
     'London', 'United Kingdom', '2017-12-21', NULL,
     'West End Production'),

    -- Sweeney Todd - Donmar 2006
    ((SELECT id FROM works WHERE slug = 'sweeney-todd'),
     'sweeney-todd-donmar-2006',
     (SELECT id FROM companies WHERE slug = 'donmar-warehouse'),
     (SELECT id FROM venues WHERE slug = 'donmar-warehouse-venue'),
     'London', 'United Kingdom', '2006-02-07', '2006-04-15',
     'Donmar Warehouse 2006'),

    -- Into the Woods - NT 2022
    ((SELECT id FROM works WHERE slug = 'into-the-woods'),
     'into-the-woods-nt-2022',
     (SELECT id FROM companies WHERE slug = 'national-theatre'),
     (SELECT id FROM venues WHERE slug = 'olivier-theatre'),
     'London', 'United Kingdom', '2022-06-23', '2022-10-08',
     'National Theatre 2022'),

    -- Phantom - Original West End 1986
    ((SELECT id FROM works WHERE slug = 'the-phantom-of-the-opera'),
     'phantom-west-end-1986',
     NULL,
     (SELECT id FROM venues WHERE slug = 'her-majestys-theatre'),
     'London', 'United Kingdom', '1986-10-09', '2023-04-22',
     'Original West End Production'),

    -- La Traviata - Met Opera 2018
    ((SELECT id FROM works WHERE slug = 'la-traviata'),
     'la-traviata-met-2018',
     (SELECT id FROM companies WHERE slug = 'metropolitan-opera'),
     (SELECT id FROM venues WHERE slug = 'metropolitan-opera-house'),
     'New York', 'United States', '2018-12-04', '2019-01-12',
     'Metropolitan Opera 2018/19 Season'),

    -- La Traviata - ROH 2019
    ((SELECT id FROM works WHERE slug = 'la-traviata'),
     'la-traviata-roh-2019',
     (SELECT id FROM companies WHERE slug = 'royal-opera-house'),
     (SELECT id FROM venues WHERE slug = 'royal-opera-house-venue'),
     'London', 'United Kingdom', '2019-01-19', '2019-02-16',
     'Royal Opera House 2019'),

    -- La Bohème - Met Opera 2020
    ((SELECT id FROM works WHERE slug = 'la-boheme'),
     'la-boheme-met-2020',
     (SELECT id FROM companies WHERE slug = 'metropolitan-opera'),
     (SELECT id FROM venues WHERE slug = 'metropolitan-opera-house'),
     'New York', 'United States', '2020-01-07', '2020-03-07',
     'Metropolitan Opera 2019/20 Season'),

    -- Streetcar - STC / Cate Blanchett 2009
    ((SELECT id FROM works WHERE slug = 'a-streetcar-named-desire'),
     'streetcar-stc-2009',
     (SELECT id FROM companies WHERE slug = 'sydney-theatre-company'),
     (SELECT id FROM venues WHERE slug = 'roslyn-packer-theatre'),
     'Sydney', 'Australia', '2009-09-05', '2009-11-07',
     'Sydney Theatre Company 2009'),

    -- King Lear - Stratford Festival 2023
    ((SELECT id FROM works WHERE slug = 'king-lear'),
     'king-lear-stratford-2023',
     (SELECT id FROM companies WHERE slug = 'stratford-festival'),
     (SELECT id FROM venues WHERE slug = 'festival-theatre-stratford'),
     'Stratford', 'Canada', '2023-06-01', '2023-10-28',
     'Stratford Festival 2023'),

    -- Hamlet - Schaubühne / Ivo van Hove 2019
    ((SELECT id FROM works WHERE slug = 'hamlet'),
     'hamlet-schaubuehne-2019',
     (SELECT id FROM companies WHERE slug = 'schaubuehne-berlin'),
     NULL,
     'Berlin', 'Germany', '2019-01-18', '2019-06-30',
     'Schaubühne 2019'),

    -- Fences - Steppenwolf 2017
    ((SELECT id FROM works WHERE slug = 'fences'),
     'fences-steppenwolf-2017',
     (SELECT id FROM companies WHERE slug = 'steppenwolf-theatre'),
     NULL,
     'Chicago', 'United States', '2017-02-02', '2017-03-26',
     'Steppenwolf 2017'),

    -- Into the Woods - Stratford 2022
    ((SELECT id FROM works WHERE slug = 'into-the-woods'),
     'into-the-woods-stratford-2022',
     (SELECT id FROM companies WHERE slug = 'stratford-festival'),
     (SELECT id FROM venues WHERE slug = 'festival-theatre-stratford'),
     'Stratford', 'Canada', '2022-05-28', '2022-10-30',
     'Stratford Festival 2022')
ON CONFLICT DO NOTHING;

-- ============================================================
-- PRODUCTION CREDITS (directors, actors)
-- ============================================================
INSERT INTO production_credits (production_id, person_id, role_type) VALUES
    -- King Lear NT 2018 — Ian McKellen
    ((SELECT id FROM productions WHERE slug = 'king-lear-nt-2018'),
     (SELECT id FROM people WHERE slug = 'ian-mckellen'), 'actor'),

    -- Hamlet Almeida 2017 — Robert Icke directed
    ((SELECT id FROM productions WHERE slug = 'hamlet-almeida-2017'),
     (SELECT id FROM people WHERE slug = 'robert-icke'), 'director'),

    -- A Doll's House Almeida 2017 — Robert Icke directed
    ((SELECT id FROM productions WHERE slug = 'dolls-house-almeida-2017'),
     (SELECT id FROM people WHERE slug = 'robert-icke'), 'director'),

    -- Angels in America NT 2017 — Marianne Elliott directed
    ((SELECT id FROM productions WHERE slug = 'angels-in-america-nt-2017'),
     (SELECT id FROM people WHERE slug = 'marianne-elliott'), 'director'),

    -- Crucible NT 2014 — Ivo van Hove directed
    ((SELECT id FROM productions WHERE slug = 'crucible-nt-2014'),
     (SELECT id FROM people WHERE slug = 'ivo-van-hove'), 'director'),

    -- Hamlet Schaubühne — Ivo van Hove directed
    ((SELECT id FROM productions WHERE slug = 'hamlet-schaubuehne-2019'),
     (SELECT id FROM people WHERE slug = 'ivo-van-hove'), 'director'),

    -- Streetcar STC 2009 — Cate Blanchett
    ((SELECT id FROM productions WHERE slug = 'streetcar-stc-2009'),
     (SELECT id FROM people WHERE slug = 'cate-blanchett'), 'actor'),

    -- Fences Roundabout 2010 — Denzel Washington
    ((SELECT id FROM productions WHERE slug = 'fences-roundabout-2010'),
     (SELECT id FROM people WHERE slug = 'denzel-washington'), 'actor'),

    -- Sweeney Todd Donmar 2006 — Sam Mendes directed
    ((SELECT id FROM productions WHERE slug = 'sweeney-todd-donmar-2006'),
     (SELECT id FROM people WHERE slug = 'sam-mendes'), 'director'),

    -- Hamilton Broadway — Lin-Manuel Miranda acted
    ((SELECT id FROM productions WHERE slug = 'hamilton-broadway-2015'),
     (SELECT id FROM people WHERE slug = 'lin-manuel-miranda'), 'actor'),

    -- Into the Woods NT 2022 — Audra McDonald
    ((SELECT id FROM productions WHERE slug = 'into-the-woods-nt-2022'),
     (SELECT id FROM people WHERE slug = 'audra-mcdonald'), 'actor')
ON CONFLICT DO NOTHING;

COMMIT;
