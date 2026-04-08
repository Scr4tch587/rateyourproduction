"use client";

import { useEffect, useState } from "react";
import { useRouter, useSearchParams } from "next/navigation";
import { WorkCard } from "@/components/work-card";
import { API_URL } from "@/lib/api";
import type { PaginatedResponse, WorkCard as WorkCardType } from "@/lib/types";
import { Input } from "@/components/ui/input";

const WORK_TYPES = ["play", "musical", "opera"] as const;
const SORT_OPTIONS = [
  { value: "top", label: "Top Rated" },
  { value: "popular", label: "Most Popular" },
  { value: "newest", label: "Newest" },
] as const;

const FIELD_LABELS = [
  ["q", "Search"],
  ["genre", "Genre"],
  ["company", "Company"],
  ["venue", "Venue"],
  ["city", "City"],
  ["country", "Country"],
  ["creator", "Creator"],
  ["person", "Person"],
  ["role_type", "Role Type"],
  ["year_from", "Year From"],
  ["year_to", "Year To"],
  ["min_rating", "Min Rating"],
  ["min_rating_count", "Min Ratings"],
] as const;

export function DiscoverContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const [works, setWorks] = useState<WorkCardType[]>([]);
  const [total, setTotal] = useState(0);
  const [loading, setLoading] = useState(true);

  const selectedTypes = searchParams.getAll("type");
  const sort = searchParams.get("sort") || "top";
  const page = Number(searchParams.get("page") || "1");

  useEffect(() => {
    const controller = new AbortController();

    async function fetchWorks() {
      setLoading(true);
      try {
        const res = await fetch(`${API_URL}/api/v1/discover?${searchParams.toString()}`, {
          credentials: "include",
          signal: controller.signal,
        });
        if (!res.ok) {
          throw new Error("discover failed");
        }
        const data: PaginatedResponse<WorkCardType> = await res.json();
        setWorks(data.data || []);
        setTotal(data.total);
      } catch {
        setWorks([]);
      } finally {
        setLoading(false);
      }
    }

    fetchWorks();
    return () => controller.abort();
  }, [searchParams]);

  function updateParams(mutator: (params: URLSearchParams) => void) {
    const params = new URLSearchParams(searchParams.toString());
    mutator(params);
    if (!params.get("sort")) params.set("sort", sort);
    if (!params.get("per_page")) params.set("per_page", "50");
    if (!params.get("page")) params.set("page", "1");
    router.push(`/discover?${params.toString()}`);
  }

  function toggleType(type: string) {
    updateParams((params) => {
      const existing = params.getAll("type");
      params.delete("type");
      if (existing.includes(type)) {
        existing.filter((value) => value !== type).forEach((value) => params.append("type", value));
      } else {
        [...existing, type].forEach((value) => params.append("type", value));
      }
      params.set("page", "1");
    });
  }

  function setTextParam(key: string, value: string) {
    updateParams((params) => {
      if (value.trim()) {
        params.set(key, value.trim());
      } else {
        params.delete(key);
      }
      if (key !== "page") {
        params.set("page", "1");
      }
    });
  }

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="border-b border-border pb-3 mb-6 flex flex-col gap-2 md:flex-row md:items-end md:justify-between">
        <div>
          <h1 className="font-serif text-2xl font-bold">Discover</h1>
          <p className="text-xs text-muted-foreground">
            Filter works by title, metadata, productions, and credited people.
          </p>
        </div>
        <select
          className="h-8 border border-input bg-background px-2 text-xs"
          value={sort}
          onChange={(e) => setTextParam("sort", e.target.value)}
        >
          {SORT_OPTIONS.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      </div>

      <div className="grid grid-cols-1 gap-8 md:grid-cols-[260px_1fr]">
        <aside className="space-y-4 rounded-sm border border-border bg-card p-4">
          <div>
            <h2 className="font-serif text-xs font-bold uppercase tracking-[0.14em] text-muted-foreground">
              Type
            </h2>
            <div className="mt-2 flex flex-wrap gap-1.5">
              {WORK_TYPES.map((type) => (
                <button
                  key={type}
                  type="button"
                  onClick={() => toggleType(type)}
                  className={`border px-2 py-1 text-[10px] uppercase tracking-wide ${
                    selectedTypes.includes(type)
                      ? "border-burgundy bg-burgundy text-primary-foreground"
                      : "border-border text-muted-foreground hover:border-foreground hover:text-foreground"
                  }`}
                >
                  {type}
                </button>
              ))}
            </div>
          </div>

          <div className="grid gap-3">
            {FIELD_LABELS.map(([key, label]) => (
              <label key={key} className="block text-[11px] text-muted-foreground">
                <span className="mb-1 block font-serif text-xs font-bold text-foreground">
                  {label}
                </span>
                <Input
                  value={searchParams.get(key) || ""}
                  onChange={(e) => setTextParam(key, e.target.value)}
                  className="h-8 bg-background text-xs"
                />
              </label>
            ))}
          </div>
        </aside>

        <section>
          <div className="mb-3 border-b border-border pb-1 text-xs text-muted-foreground">
            {loading ? "Searching..." : `${total} result${total === 1 ? "" : "s"}`}
          </div>

          {works.length > 0 ? (
            <div className="rounded-sm border border-border bg-card px-3">
              {works.map((work) => (
                <WorkCard key={work.id} work={work} />
              ))}
            </div>
          ) : (
            !loading && (
              <div className="rounded-sm border border-border bg-card px-4 py-8 text-center text-xs text-muted-foreground italic">
                No works found. Adjust the archival filters and try again.
              </div>
            )
          )}

          {total > 50 && (
            <div className="mt-4 flex items-center justify-center gap-4 text-xs">
              {page > 1 && (
                <button type="button" onClick={() => setTextParam("page", String(page - 1))}>
                  Previous
                </button>
              )}
              <span className="text-muted-foreground">
                Page {page} of {Math.ceil(total / 50)}
              </span>
              {page < Math.ceil(total / 50) && (
                <button type="button" onClick={() => setTextParam("page", String(page + 1))}>
                  Next
                </button>
              )}
            </div>
          )}
        </section>
      </div>
    </div>
  );
}
