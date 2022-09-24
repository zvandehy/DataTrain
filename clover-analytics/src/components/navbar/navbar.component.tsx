import { MilitaryTechRounded } from "@mui/icons-material";
import { Icon, Toolbar, Typography } from "@mui/material";
import "./navbar.component.css";

interface NavBarProps {}

const NavBar: React.FC<NavBarProps> = ({}: NavBarProps) => {
  return (
    <Toolbar sx={{ justifyContent: "center" }}>
      <Icon sx={{ justifySelf: "start" }}>
        <MilitaryTechRounded />
      </Icon>
      <Typography>Child 1</Typography>
      <Typography>Child 2</Typography>
      <Typography>Child 3</Typography>
    </Toolbar>
  );
};

export default NavBar;
