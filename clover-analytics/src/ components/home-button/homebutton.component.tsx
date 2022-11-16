import FilterVintage from "@mui/icons-material/FilterVintage";
import { Button, Typography } from "@mui/material";

export const HomeButton = (props: any) => {
  return (
    <Button
      sx={{
        color: "lightgreen",
        background: "transparent",
        "&:hover": {
          backgroundColor: "green",
          color: "white",
        },
      }}
    >
      <FilterVintage sx={{ mr: { xs: 0, sm: 2 } }} />
      <Typography
        variant="h6"
        noWrap
        component="div"
        sx={{ display: { xs: "none", sm: "block" } }}
      >
        CLOVER ANALYTICS
      </Typography>
    </Button>
  );
};
