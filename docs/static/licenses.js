const readJson = (url, callback) => {
    const ajax = new XMLHttpRequest()
    ajax.open('GET', `${window.location.protocol}//${window.location.host}/` + url)
    ajax.send()
    ajax.addEventListener("load", () => {
        callback(JSON.parse(ajax.response))
    })
}

const buildElement = (tagName, attributes, textNode) => {
    const element = document.createElement(tagName)
    if (attributes != undefined) {
        Object.keys(attributes).forEach(key => element.setAttribute(key, attributes[key]))
    }
    if (textNode != undefined) {
        element.appendChild(document.createTextNode(textNode))
    }
    return element
}

const removeAllChildren = (element) => {
    while (element.firstChild) {
        element.removeChild(element.firstChild)
    }
}

const updateInfo = (json) => {
    const timestampElement = document.getElementById('timestamp')
    timestampElement.appendChild(document.createTextNode(json.timestamp))
    const commitElement = document.getElementById('commitId')
    const anchor = buildElement("a", { "href": "https://github.com/spdx/license-list-XML/commit/" + json["git-commit-id"] })
    anchor.appendChild(document.createTextNode(json["git-commit-id"].substring(0, 8)))
    commitElement.appendChild(anchor)
}
const layoutLicenses = (json, condition) => {
    const buildLicenseInfo = (license) => {
        const td = buildElement("td", {})
        if (license["osi-approved"]) {
            td.appendChild(buildElement("img", { "src": "../images/approved.svg", "alt": "OSI Approved" }))
        }
        if (license["deprecated"]) {
            td.appendChild(buildElement("img", { "src": "../images/deprecated.svg", "alt": "Deprecated" }))
        }
        return td
    }
    const buildUrls = (urls) => {
        const td = buildElement('td')
        urls.forEach((url, index) => {
            const anchor = buildElement("a", { "href": url })
            const img = buildElement("img", { "width": 16, "height": 16, "src": "../images/external-link.svg", "title": url })
            anchor.appendChild(img)
            td.appendChild(anchor)
        })
        return td
    }
    const buildName = (name) => {
        const ruby = buildElement("div", { "class": "license-name" })
        const fullName = buildElement("span", { "class": "license-full-name" }, name.full)
        const shortName = buildElement("span", { "class": "license-short-name" }, name.short)
        ruby.appendChild(fullName)
        ruby.appendChild(shortName)
        const td = buildElement("td", {})
        td.appendChild(ruby)
        return td
    }
    const buildLicenseRow = (license, index) => {
        const tr = buildElement('tr')
        tr.appendChild(buildElement('td', {}, index + 1))
        tr.appendChild(buildName(license.name))
        tr.appendChild(buildLicenseInfo(license))
        tr.appendChild(buildUrls(license.urls))
        return tr
    }
    const showLicense = (license) => {
        if (license["osi-approved"] && license["deprecated"]) {
            return condition["osi-deprecated"]
        } else if (license["osi-approved"] && !license["deprecated"]) {
            return condition["osi"]
        } else if (!license["osi-approved"] && license["deprecated"]) {
            return condition["deprecated"]
        }
        return condition["non-osi"]
    }
    const filterByName = (license) => {
        if (condition.filterText == "") {
            return true
        }
        const fullName = license.name.full.toLowerCase()
        const shortName = license.name.short.toLowerCase()
        return shortName.indexOf(condition.filterText) >= 0 ||
            fullName.indexOf(condition.filterText) >= 0
    }
    const listLicenses = (body, licenses) => {
        const results = licenses.filter(showLicense)
            .filter(filterByName)
            .map((license, index) => buildLicenseRow(license, index))
        results.forEach(node => body.appendChild(node))
        return results.length
    }
    const updateCondition = (c) => {
        const array = ["osi", "osi-deprecated", "non-osi", "deprecated"]
        array.forEach(key => {
            const element = document.getElementById("checkbox-" + key)
            c[key] = element.checked
        })
        const name = document.getElementById("filter-text").value
        c.filterText = name.toLowerCase()
    }
    const updateCount = (count) => {
        const element = document.getElementById("license-count")
        removeAllChildren(element)
        if (count == 1) {
            element.appendChild(document.createTextNode(`${count} license`))
        } else {
            element.appendChild(document.createTextNode(`${count} licenses`))
        }
    }
    updateCondition(condition)
    const base = document.getElementById('licenses')
    const tbody = base.getElementsByTagName('tbody')[0]
    removeAllChildren(tbody)
    const count = listLicenses(tbody, json.licenses)
    updateCount(count)
}

const initLicenses = (callback) => {
    readJson("lioss/spdx_licenses.json", (json) => {
        updateInfo(json)
        layoutLicenses(json, { "osi": true, "osi-deprecated": true, "non-osi": true, "deprecated": true })
        callback(json)
    })
}
