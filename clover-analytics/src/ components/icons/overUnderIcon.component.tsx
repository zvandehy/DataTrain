import { DoubleArrow } from "@material-ui/icons";
import UnfoldLessDoubleIcon from "@mui/icons-material/UnfoldLessDouble";
import { Icon, Typography } from "@mui/material";
import { makeStyles } from "@mui/styles";
import { COLORS } from "../../shared/styles/constants";

const overUnderStyle = makeStyles({
  over: {
    color: COLORS.HIGHER_DARK,
    bgColor: "white",
    transform: "rotate(270deg)",
  },
  under: {
    color: COLORS.LOWER_DARK,
    bgColor: "white",
    transform: "rotate(90deg)",
  },
  push: {
    color: COLORS.PUSH_DARK,
    transform: "rotate(90deg)",
  },
});

interface OverUnderIconProps {
  overUnder: string;
  size?: number;
}
// TODO: Return an icon instead of a react element
export const OverUnderIcon: React.FC<OverUnderIconProps> = ({
  overUnder,
  size,
}: OverUnderIconProps) => {
  const classes = overUnderStyle();
  return (
    <Icon>
      {overUnder.toUpperCase() === "OVER" ? (
        <DoubleArrow className={classes.over} style={{ lineHeight: 0.8 }} />
      ) : overUnder.toUpperCase() === "UNDER" ? (
        <DoubleArrow className={classes.under} />
      ) : overUnder.toUpperCase() === "PUSH" ? (
        <UnfoldLessDoubleIcon className={classes.push} />
      ) : (
        <></>
      )}
    </Icon>
  );
};
// TODO: Return an icon instead of a react element
export const OverUnderTypography: React.FC<OverUnderIconProps> = ({
  overUnder,
  size,
}: OverUnderIconProps) => {
  return (
    <Typography
      //   variant={"overline"}
      letterSpacing={1.5}
      fontSize={size ?? 12}
      fontWeight={600}
      textTransform={"uppercase"}
    >
      {overUnder}
    </Typography>
  );
};
