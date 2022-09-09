import styled from "@emotion/styled";
import { TableCell, tableCellClasses, tableHeadClasses } from "@mui/material";

export const StyledTableCell = styled(TableCell)(({ theme }) => ({
  [`&.${tableCellClasses.head}`]: {
    backgroundColor: "inherit",
    padding: "3px",
    paddingBottom: "0px",
    color: "inherit",
    fontSize: "1rem",
    alignItems: "center",
    textAlign: "center",
    fontWeight: "bold",
  },
  [`&.${tableHeadClasses.root}`]: {
    backgroundColor: "inherit",
    padding: "3px",
    paddingBottom: "0px",
    color: "inherit",
    fontSize: "1rem",
    alignItems: "center",
    textAlign: "center",
    fontWeight: "bold",
  },
  [`&.${tableCellClasses.body}`]: {
    padding: "3px",
    fontSize: "1rem",
    color: "inherit",
    alignItems: "center",
    textAlign: "center",
  },
}));
