import styled from "@emotion/styled";
import { TableCell, tableCellClasses } from "@mui/material";

export const StyledTableCell = styled(TableCell)(({ theme }) => ({
  [`&.${tableCellClasses.head}`]: {
    backgroundColor: "inherit",
    // padding: "15px min(1.5vw, 1.5rem);",
    padding: "5px",
    paddingBottom: "0px",
    color: "inherit",
    fontSize: "unset",
    alignItems: "center",
    textAlign: "center",
    fontWeight: "bold",
  },
  [`&.${tableCellClasses.body}`]: {
    // padding: "min(1vw, 1.5rem);",
    padding: "5px",
    fontSize: "1rem",
    color: "inherit",
    alignItems: "center",
    textAlign: "center",
  },
}));
