async function fetchData(url, options) {
    const response = await fetch(url, options);
    const data = await response.json();
    return data;
}
const root = "https://localhost"
const case_id = await fetchData(root + "/Selection/Case/Create", {
    "body": "{\"client_id\":\"C123\",\"description\":\"Gebroken \",\"name\":\"Andrea\"}",
    "method": "POST"
});

const techs = await fetchData(root + "/Selection/Tech")
const tech_id = techs[0].id

await fetch(root + `/Selection/Tech/${tech_id}/Shortlist`, {
    "body": "case.id=" + case_id,
    "method": "POST",
});

const techChoices = await fetchData(root + "/Selection/TechChoice/" + case_id)
const techChoice_id = techChoices[0].id

await fetch(`https://localhost/Selection/TechChoice/${techChoice_id}/Choose`, {
    "body": `case.id=${case_id}&reasoning=TEST`,
    "method": "POST"
});