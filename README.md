# TeamFeedback

[![License](https://img.shields.io/badge/license-BSD--3-green)](LICENSE)
[![Contributor Covenant](https://img.shields.io/badge/Contributor%20Covenant-2.1-4baaaa.svg)](CODE_OF_CONDUCT.md)
![Go CI](https://github.com/cprbucat2/team-feedback/actions/workflows/go-ci.yml/badge.svg)

TeamFeedback is a Dockerized feedback solution for teams, with both a web server and serverless solution both written in Go. We aim to ease the submission and redistribution of course feedback, including any aggregation, anonymization, or editing that the coordinators wish to perform. Coordinators can create teams of any size and view/vet all data. Team members can only view feedback as it is redistributed to them (e.g., after aggregation). The web application has different views for each role, allowing coordinators to manage teams and review data while students submit and review feedback. The serverless TeamFeedback exists to support coordinators who do not want to host the Dockerized database and web server. The local tool can generate static pages (with local resources) that collect, aggregate, and distribute data to the specific team members. Our stack includes Docker, MySQL, Go Gin, HTML, CSS, and JavaScript.
