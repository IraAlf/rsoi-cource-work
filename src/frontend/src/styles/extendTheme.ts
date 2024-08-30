import { extendTheme } from "@chakra-ui/react";

const theme = extendTheme({
  colors: {
    primary: "#4A6572", // Deep blue-grey
    "primary-dark": "#344955", // Darker blue-grey
    secondary: "#F0F4F8", // Light grey
    third: "#E2E8F0", // Slightly darker grey

    "bg": "#F8FAFC", // Very light grey
    "bg-alt": "#E5E7EB", // Lighter grey
    "text": "#2D3748", // Dark grey
    "title": "#1A202C", // Darkest grey
    "button": "#3182CE", // Blue
    "accent-1": "#718096", // Light blue-grey
    "accent-2": "#0EA5E9", // Blue
    "accent-3": "#4299E1", // Light blue
    "add-1": "#A0AEC0", // Light grey
    "add-2": "#607D8B", // Dark grey
  },
  fonts: {
    heading: "Poppins",
    body: "Roboto",
  },

  textStyles : {
    headline: {
      fontFamily: 'Poppins',
      fontStyle: 'normal',
      fontWeight: 'bold',
      fontSize: '66px',
      lineHeight: '110%',
    },

    header: {
      fontFamily: 'Poppins',
      fontStyle: 'normal',
      fontWeight: 'bold',
      fontSize: '41px',
      lineHeight: '140%',
    },

    subtitle: {
      fontFamily: 'Poppins',
      fontStyle: 'normal',
      fontWeight: 'bold',
      fontSize: '26px',
      lineHeight: '140%',
    },

    body: {
      fontFamily: 'Roboto',
      fontStyle: 'normal',
      fontWeight: 'normal',
      fontSize: '16px',
      lineHeight: '200%',
    },

    'alt-body': {
      fontFamily: 'Roboto',
      fontStyle: 'normal',
      fontWeight: 'normal',
      fontSize: '16px',
      lineHeight: '20px',
      letterSpacing: '-0.045em'
    },

    caption: {
      fontFamily: 'Roboto',
      fontStyle: 'normal',
      fontWeight: 'normal',
      fontSize: '12px',
      lineHeight: '200%',
      textTransform: 'uppercase',
    },
  }
});

export default theme;