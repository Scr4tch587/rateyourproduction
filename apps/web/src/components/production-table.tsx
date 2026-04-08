import Link from "next/link";
import { Rating } from "@/components/rating";
import {
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableHeader,
  TableRow,
} from "@/components/ui/table";
import type { Production } from "@/lib/types";

interface ProductionTableProps {
  productions: Production[];
}

export function ProductionTable({ productions }: ProductionTableProps) {
  if (productions.length === 0) {
    return (
      <div className="text-xs text-muted-foreground py-4 italic">
        No productions on record.
      </div>
    );
  }

  return (
    <Table>
      <TableHeader>
        <TableRow className="text-[10px] uppercase tracking-wide text-muted-foreground">
          <TableHead className="font-sans font-semibold">Production</TableHead>
          <TableHead className="font-sans font-semibold">Company</TableHead>
          <TableHead className="font-sans font-semibold">Venue</TableHead>
          <TableHead className="font-sans font-semibold">Year</TableHead>
          <TableHead className="font-sans font-semibold text-right">Rating</TableHead>
          <TableHead className="font-sans font-semibold text-right">Ratings</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {productions.map((p, i) => (
          <TableRow
            key={p.id}
            className={`text-xs ${i % 2 === 1 ? "bg-muted/40" : ""}`}
          >
            <TableCell>
              <Link
                href={`/productions/${p.id}`}
                className="hover:text-burgundy transition-colors"
              >
                {p.production_label || p.city || p.slug}
              </Link>
            </TableCell>
            <TableCell className="text-muted-foreground">
              {p.company_name || "\u2014"}
            </TableCell>
            <TableCell className="text-muted-foreground">
              {p.venue_name || "\u2014"}
            </TableCell>
            <TableCell className="text-muted-foreground">
              {p.start_date ? new Date(p.start_date).getFullYear() : "\u2014"}
            </TableCell>
            <TableCell className="text-right">
              {p.rating_count > 0 ? (
                <Rating rating={p.average_rating} count={0} />
              ) : (
                <span className="text-muted-foreground italic">\u2014</span>
              )}
            </TableCell>
            <TableCell className="text-right text-muted-foreground">
              {p.rating_count > 0 ? p.rating_count : "\u2014"}
            </TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
