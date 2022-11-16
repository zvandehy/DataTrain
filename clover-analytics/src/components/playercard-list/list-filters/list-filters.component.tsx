import { TextField } from "@mui/material";
import { BETTING_CATEGORIES } from "../../../shared/constants";
import { Stat } from "../../../shared/interfaces/stat.interface";
import { Option } from "../../../shared/interfaces/option.interface";
import AutocompleteFilter from "../../autocomplete-filter/autocomplete-filter.component";
import "./list-filters.component.css";

// List filters are used to filter a list of projections that have already been loaded.

interface PlayerListFiltersProps {
  onSearchChange: (value: string) => void;
  onSortSelect: (value: string) => void;
  onStatSelect: (value: Stat) => void;
}

let StatOptions: Option<Stat | undefined>[] = BETTING_CATEGORIES.map((stat) => {
  return {
    id: stat,
    label: stat.display,
  };
});
StatOptions.unshift({
  id: undefined,
  label: "ANY",
});

const PlayerListFilters: React.FC<PlayerListFiltersProps> = ({
  onSearchChange,
  onSortSelect,
  onStatSelect,
}: PlayerListFiltersProps) => {
  return (
    <div id="list-filters" className={"filters-wrapper"}>
      <AutocompleteFilter
        options={StatOptions}
        onChange={onStatSelect}
        label="Stat"
      />
      <AutocompleteFilter
        options={[
          { label: "Confidence", id: "confidence" },
          { label: "Name", id: "name" },
        ]}
        onChange={onSortSelect}
        label="Sort by"
      />
      <TextField
        label="Search..."
        onChange={(e) => onSearchChange(e.target.value)}
        sx={
          {
            // textAlign: "end",
            // width: "fit-content",
          }
        }
      />
    </div>
  );
};

export default PlayerListFilters;
