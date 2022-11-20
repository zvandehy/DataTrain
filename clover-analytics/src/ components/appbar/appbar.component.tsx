import SearchIcon from "@mui/icons-material/Search";
import { Icon, IconButton } from "@mui/material";
import { AdapterMoment } from "@mui/x-date-pickers/AdapterMoment";
import AppBar from "@mui/material/AppBar";
import Box from "@mui/material/Box";
import InputBase from "@mui/material/InputBase";
import Menu from "@mui/material/Menu";
import MenuItem from "@mui/material/MenuItem";
import { alpha, styled } from "@mui/material/styles";
import Toolbar from "@mui/material/Toolbar";
import { DatePicker, LocalizationProvider } from "@mui/x-date-pickers";
import moment from "moment";
import * as React from "react";
import { Calendar } from "react-calendar";
import { HomeButton } from "../home-button/homebutton.component";
import { ModelButton } from "../styled-model-button/model-button.component";

const Search = styled("div")(({ theme }) => ({
  position: "relative",
  borderRadius: theme.shape.borderRadius,
  backgroundColor: alpha(theme.palette.common.white, 0.15),
  "&:hover": {
    backgroundColor: alpha(theme.palette.common.white, 0.25),
  },
  marginRight: theme.spacing(2),
  marginLeft: 0,
  width: "100%",
  [theme.breakpoints.up("xs")]: {
    marginLeft: theme.spacing(3),
    width: "auto",
  },
}));

const SearchIconWrapper = styled("div")(({ theme }) => ({
  padding: theme.spacing(0, 2),
  height: "100%",
  position: "absolute",
  pointerEvents: "none",
  display: "flex",
  alignItems: "center",
  justifyContent: "center",
}));

const StyledInputBase = styled(InputBase)(({ theme }) => ({
  color: "inherit",
  "& .MuiInputBase-input": {
    padding: theme.spacing(1, 1, 1, 0),
    // vertical padding + font size from searchIcon
    paddingLeft: `calc(1em + ${theme.spacing(4)})`,
    transition: theme.transitions.create("width"),
    width: "100%",
    [theme.breakpoints.up("md")]: {
      width: "20ch",
    },
  },
}));

export interface AppBarProps {
  onDateSelect: (date: string) => void;
  date: string;
}

export const PrimarySearchAppBar = ({ onDateSelect, date }: AppBarProps) => {
  const [anchorEl, setAnchorEl] = React.useState<null | HTMLElement>(null);
  const [calendarOpen, setCalendarOpen] = React.useState(false);

  const openCalendar = (open: boolean) => {
    setCalendarOpen(open);
    console.log("openCalendar", open);
  };

  const isMenuOpen = Boolean(anchorEl);

  const handleProfileMenuOpen = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };

  const menuId = "primary-search-account-menu";
  const renderMenu = (
    <Menu
      anchorEl={anchorEl}
      anchorOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      id={menuId}
      keepMounted
      transformOrigin={{
        vertical: "top",
        horizontal: "right",
      }}
      open={isMenuOpen}
      onClose={handleMenuClose}
    >
      <MenuItem onClick={handleMenuClose}>Profile</MenuItem>
      <MenuItem onClick={handleMenuClose}>My account</MenuItem>
    </Menu>
  );

  return (
    <Box sx={{ flexGrow: 1 }}>
      <AppBar position="static">
        <Toolbar>
          <HomeButton />
          <Search>
            <SearchIconWrapper>
              <SearchIcon />
            </SearchIconWrapper>
            <StyledInputBase
              placeholder="Search Player…"
              inputProps={{ "aria-label": "search" }}
              sx={{ width: "100%" }}
            />
          </Search>
          <LocalizationProvider dateAdapter={AdapterMoment}>
            <DatePicker
              value={date}
              open={calendarOpen}
              onOpen={() => openCalendar(true)}
              onClose={() => openCalendar(false)}
              onChange={(newValue) =>
                onDateSelect(moment(newValue).format("YYYY-MM-DD") || "")
              }
              renderInput={(props) => (
                <IconButton onClick={() => openCalendar(true)}>
                  <Icon>today</Icon>
                </IconButton>
              )}
            />
          </LocalizationProvider>
          <Box sx={{ flexGrow: 5 }} />
          <Box sx={{ flexGrow: 1 }}>
            <ModelButton />
          </Box>
          {/* <Box>
            <IconButton
              size="large"
              edge="end"
              aria-label="account of current user"
              aria-controls={menuId}
              aria-haspopup="true"
              onClick={handleProfileMenuOpen}
              color="inherit"
            >
              <AccountCircle />
            </IconButton>
          </Box> */}
        </Toolbar>
      </AppBar>
      {renderMenu}
    </Box>
  );
};
