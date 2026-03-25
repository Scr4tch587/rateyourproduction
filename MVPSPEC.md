# RateYourProduction MVP Spec

**Source of Truth v1 (Locked)**

This document defines the complete MVP: product, data model, backend, discovery system, ranking, and UI.
If something is not in this spec, it is **out of scope**.

---

# 1. Product Definition

## 1.1 Product Summary

RateYourProduction is a web platform for logging, rating, reviewing, and discovering **theatre works** and their **productions**.

Supported domains:

* plays
* musicals
* operas

## 1.2 Core Concept

The platform is built on a **title-first system**:

* **Work (canonical)** → title-level object
* **Production** → specific staging of a work
* **Log** → user interaction (attendance + rating + review)

## 1.3 Core User Actions

Users can:

* search for works
* browse productions under a work
* log what they’ve seen
* rate and review
* explore works via advanced filters
* view user profiles and logs
* submit missing productions

---

# 2. Core Product Decisions (Non-Negotiable)

## 2.1 Canonical Hierarchy

```
Work → Production → Log
```

* All discovery is anchored on **Work**
* Productions exist only within a Work
* Logs attach to:

  * production (preferred)
  * work (fallback)

## 2.2 Discovery Model

* Discovery returns **works**, not productions
* Filters operate across:

  * work metadata
  * production metadata
  * production credits (people)

## 2.3 Logging Model

* Logging is the **primary action**
* Reviews are optional
* Production selection is optional but encouraged

---

# 3. MVP Scope

## 3.1 Included

* Auth
* Work pages
* Production pages
* Logging (rating + review)
* User profiles
* Advanced discovery (ALL filter tiers)
* Weighted ranking system
* User-submitted productions
* Basic admin panel

## 3.2 Excluded

* follows/friends
* comments
* likes
* lists
* recommendations
* notifications
* mobile app
* full moderation system
* cast completeness requirements

---

# 4. Data Model

## 4.1 Users

* id
* username
* display_name
* avatar_url
* created_at

---

## 4.2 Works

* id
* slug
* title
* normalized_title
* type (`play`, `musical`, `opera`)
* description
* premiere_year
* average_rating
* rating_count
* weighted_score
* created_at
* updated_at

---

## 4.3 WorkCreators

* work_id
* person_id
* role_type (`playwright`, `composer`, etc.)

---

## 4.4 Genres

* id
* name
* slug

---

## 4.5 WorkGenres

* work_id
* genre_id

---

## 4.6 Productions

* id
* work_id
* slug
* company_id
* venue_id (nullable)
* city
* country
* start_date
* end_date
* production_label (optional)
* average_rating
* rating_count
* weighted_score
* created_at

---

## 4.7 Companies

* id
* name
* slug
* city
* country

---

## 4.8 Venues

* id
* name
* slug
* city
* country

---

## 4.9 People

* id
* name
* slug

---

## 4.10 ProductionCredits

* production_id
* person_id
* role_type (`actor`, `director`, etc.)

---

## 4.11 Logs

* id
* user_id
* work_id
* production_id (nullable)
* seen_date
* rating (0.5–5.0)
* review_text (nullable)
* created_at

---

# 5. Discovery System

## 5.1 Entry Point

Dedicated page: `/discover`

## 5.2 Filter System (ALL INCLUDED)

### Tier 1

* type
* genre
* company
* venue
* city
* country
* year range
* min rating
* min rating count

### Tier 2

* creator (playwright/composer)
* language (optional)
* production tags (optional)

### Tier 3 (REQUIRED)

* person (actor/director/etc.)
* role type filter

---

## 5.3 Filter Logic

* Within same category → OR
* Across categories → AND

Example:

* genre = tragedy OR drama
  AND
* company = Stratford
  AND
* year = 2018–2024

---

## 5.4 Result Type

Returns **Work cards**

Each result includes:

* title
* type
* genres
* average rating
* number of ratings
* weighted score
* production count
* matched metadata (optional highlight)

---

# 6. Ranking Algorithm

## 6.1 Formula (Required)

[
score = \frac{v}{v+m}R + \frac{m}{v+m}C
]

Where:

* `R` = average rating
* `v` = number of ratings
* `C` = global average rating
* `m` = constant (default: 10)

---

## 6.2 Stored Values

Both stored and updated:

### Work

* average_rating
* rating_count
* weighted_score

### Production

* same fields

---

## 6.3 Sorting Modes

* Top Rated (weighted score)
* Most Popular (rating_count)
* Newest (production date)

---

# 7. User Flows

## 7.1 Log Flow

1. Search work
2. Open work page
3. Select production OR skip
4. Input:

   * date
   * rating
   * review (optional)
5. Submit

---

## 7.2 Discovery Flow

1. Open `/discover`
2. Apply filters
3. Browse work results
4. Open work page
5. Select production
6. Log

---

## 7.3 Missing Production

User can:

* add production under a work

Fields:

* company
* venue
* city
* date

---

# 8. UI / Design System (RYM-Inspired)

## 8.1 Core Principle

> Dense, information-first UI inspired by RateYourMusic, modernized for usability

---

## 8.2 Theme

* Dark mode default
* High contrast text
* Minimal gradients
* No large empty spaces

---

## 8.3 Layout Rules

* Grid-based
* Table-heavy
* Data-first
* No hero sections
* Minimal imagery reliance

---

## 8.4 Typography

* compact text
* small base font
* tight spacing
* emphasis via weight, not size

---

## 8.5 Components

## Work Card

* title
* type
* genres
* rating (stars + number)
* rating count
* production count

---

## Rating Display

Always show:

```
★★★★☆ 4.2 (128)
```

Used everywhere.

---

## Work Page Layout

### Left Panel

* title
* type
* creators
* genres
* rating
* stats

### Right Panel

* production list
* recent logs

---

## Production List

Table format:

Columns:

* production name
* company
* venue
* year
* rating
* rating count

---

## Discovery Page

Left sidebar:

* filters

Main area:

* results list

Top bar:

* sort dropdown

---

## Profile Page

* username
* stats
* log list

Log entry format:

```
Hamlet — Stratford 2023
★★★★☆
“Review text…”
```

---

## Design Constraints

* no large cards
* no infinite whitespace
* no overly modern "startup UI"
* prioritize scannability over aesthetics

---

# 9. Backend Architecture

## Stack

* Go API
* Postgres (Supabase)
* Redis
* Next.js frontend

---

## Required Systems

* REST API
* Event-style logging
* Derived rating updates
* Background worker (for aggregates)

---

# 10. Deployment

* Frontend → Vercel
* Backend → Railway
* DB/Auth/Storage → Supabase
* Cache → Redis

---

# 11. Observability

* Sentry required

Track:

* errors
* slow endpoints
* crashes

---

# 12. Admin (Minimal)

Internal only:

* create/edit works
* create/edit productions
* delete duplicates
* manage submissions

---

# 13. Launch Requirements

Before launch:

* seeded works (hundreds+)
* seeded productions
* basic credit data (for filters)
* working discovery system
* working logging flow

---

# Final Definition

## MVP is complete when:

A user can:

* search for a work
* browse its productions
* log what they saw
* rate and review
* discover works using filters across:

  * genre
  * company
  * time
  * people
* see results ranked by weighted community score

---

# Guiding Principle

> This is not a theatre app with filters.
> This is a structured discovery + ranking engine with theatre as the domain.

---

If you build exactly this — no more, no less — you will have:

* a real product
* real users
* strong backend signal
* a very defensible and interesting project

Next step (if you want):
I can turn this into a **2-week execution plan with exact order of implementation**.
