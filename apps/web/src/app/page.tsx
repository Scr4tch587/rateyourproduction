import Link from "next/link";

export default function Home() {
  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <section className="mb-8">
        <h1 className="text-lg font-bold mb-1">RateYourProduction</h1>
        <p className="text-xs text-muted-foreground">
          Log, rate, and discover plays, musicals, and operas.
        </p>
      </section>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <section>
          <h2 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-3">
            Top Rated Works
          </h2>
          <div className="text-xs text-muted-foreground">
            No works yet.{" "}
            <Link href="/discover" className="underline hover:text-foreground">
              Discover works
            </Link>
          </div>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-3">
            Recent Logs
          </h2>
          <div className="text-xs text-muted-foreground">No logs yet.</div>
        </section>

        <section>
          <h2 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-3">
            Popular Productions
          </h2>
          <div className="text-xs text-muted-foreground">
            No productions yet.
          </div>
        </section>
      </div>
    </div>
  );
}
