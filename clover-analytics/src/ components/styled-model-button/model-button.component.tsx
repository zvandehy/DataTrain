import { Button } from "@mui/material";
import { styled } from "@mui/material/styles";

export const ModelButtonStyle = styled(Button)(({ theme }) => ({
  color: "white",
  background:
    "linear-gradient(205deg, rgba(255,175,29,1) 0%, rgba(252,70,107,1) 100%)",
  "&:hover": {
    backgroundColor: theme.palette.primary.dark,
  },
  width: "100%",
  maxWidth: "300px",
}));

export const ModelButton = (props: any) => {
  return <ModelButtonStyle {...props}>CUSTOM MODEL</ModelButtonStyle>;
};
