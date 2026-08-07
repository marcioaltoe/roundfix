## {{ version.tag }} - {{ date | date(format="%Y-%m-%d") }}
{% for type, typed_commits in commits | sort(attribute="type") | group_by(attribute="type") -%}
#### {{ type | upper_first }}
{% for commit in typed_commits -%}
{% if commit.scope -%}
- **({{ commit.scope }})** {{ commit.summary }}
{% else -%}
- {{ commit.summary }}
{% endif -%}
{% endfor -%}
{% endfor -%}
