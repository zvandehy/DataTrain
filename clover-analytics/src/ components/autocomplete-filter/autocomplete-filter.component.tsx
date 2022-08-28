import { Autocomplete, TextField } from "@mui/material";
import { Option } from "../../shared/interfaces/option.interface";
import "./autocomplete-filter.component.css";

interface AutocompleteFilterProps {
  options: Option<any>[];
  onChange: Function;
  label: string;
  width?: number;
}

const AutocompleteFilter: React.FC<AutocompleteFilterProps> = ({
  options,
  onChange,
  label,
  width,
}: AutocompleteFilterProps) => {
  return (
    <Autocomplete
      id="filter-select"
      className={"filter-select"}
      options={options}
      ListboxProps={{
        style: {
          color: "black",
        },
      }}
      renderInput={(params) => <TextField {...params} label={label} />}
      onChange={(event, value) => {
        onChange(value?.id || options[0].id);
      }}
      style={{ width: `${width}px`, color: "black" }}
    />
  );
};

export default AutocompleteFilter;
