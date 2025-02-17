// src/access.ts
export default function access(initialState: { currentUser?: API.CurrentUser | undefined }) {
  const { currentUser } = initialState || {};
  console.info(currentUser)
  // return {
  //    canAdmin: currentUser && currentUser.data.admin == true,
  //  };
  if (currentUser == undefined) {
    return {
      canAdmin: false,
    };
  } else {
    return {
      canAdmin: currentUser.admin && currentUser.admin == true,
    };
  }

}
