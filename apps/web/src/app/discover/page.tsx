export default function DiscoverPage() {
  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <h1 className="text-lg font-bold mb-4">Discover</h1>
      <div className="grid grid-cols-1 md:grid-cols-[240px_1fr] gap-6">
        <aside className="space-y-4">
          <div>
            <h2 className="text-xs font-semibold uppercase tracking-wide text-muted-foreground mb-2">
              Filters
            </h2>
            <div className="space-y-2 text-xs text-muted-foreground">
              <p>Type, Genre, Company, Venue, City, Country, Year range, Rating</p>
              <p>Creator, Person, Role type</p>
            </div>
          </div>
        </aside>
        <section>
          <div className="flex items-center justify-between mb-4">
            <span className="text-xs text-muted-foreground">0 results</span>
            <select className="text-xs bg-card border border-border rounded px-2 py-1">
              <option>Top Rated</option>
              <option>Most Popular</option>
              <option>Newest</option>
            </select>
          </div>
          <div className="text-xs text-muted-foreground">
            Apply filters to discover works.
          </div>
        </section>
      </div>
    </div>
  );
}
